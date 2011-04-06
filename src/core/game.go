// Copyright 2010, 2011 The ghack Authors. All rights reserved.
// Use of this source code is governed by the GNU General Public License
// version 3 (or any later version). See the file COPYING for details.

package core

import (
    "reflect"
    "time"
)

var tick_rate int64 = 60            // Ticks per second
var skip_ns int64 = 1e9 / tick_rate // Nanosecond interval per tick

// Manages game data and runs the main loop.
type Game struct {
    svc     ServiceContext
    ents    map[chan Msg]Entity
    nextUid UniqueId
    // This function is used to spawn new players when a client joins. As this
    // service has no concept of any specific entities it must be set by the specific game.
    PlayerFunc func(UniqueId) Entity
}

func NewGame(svc ServiceContext) *Game {
    var uid UniqueId = 0 // Game uid is always zero
    ents := make(map[chan Msg]Entity)
    return &Game{svc, ents, uid + 1, nil}
}

func (g *Game) Run(input chan Msg) {
    g.waitOnServiceStart(input)

    // List of up to date entities
    updated := make(map[chan Msg]bool, len(g.ents))
    tick_msg := MsgTick{input}

    for {
        tick_start := time.Nanoseconds()

        // Tell all the entities that a new tick has started
        ent_num := len(g.ents) // Store ent count for *this* tick
        for ent := range g.ents {
            ent <- tick_msg
        }

        // Listen for any service messages
        // Break out of loop once all entities have updated
        for {
            msg := <-input
            switch m := msg.(type) {
            case MsgTick:
                updated[m.Origin] = true // bool value doesn't matter
                if len(updated) == ent_num {
                    // Clear out list, use current number of entities for next tick
                    updated = make(map[chan Msg]bool, len(g.ents))
                    goto update_end
                }
            case MsgListEntities:
                m.Reply <- g.makeEntityList()
            case MsgSpawnPlayer:
                g.spawnPlayer(m)
            }
        }
    update_end:
        g.svc.Comm <- tick_msg

        sleep_ns := (tick_start + skip_ns) - time.Nanoseconds()
        time.Sleep(sleep_ns)
    }
}

func (g *Game) AddEntity(ent Entity, ch chan Msg) {
    g.ents[ch] = ent
}

func (g *Game) makeEntityList() MsgListEntities {
    list := make([]*EntityDesc, len(g.ents))

    i := 0
    for ch, ent := range g.ents { // Fill the len(g.ents) slots in list
        list[i] = NewEntityDesc(ent, ch)
        i++
    }
    return MsgListEntities{nil, list}
}

// Returns the next available unique id
func (g *Game) GetUid() UniqueId {
    uid := g.nextUid
    g.nextUid++
    return uid
}

// Returns once all services have signalled that they have started.
// Automatically accounts for a variable number of services as contained in the
// ServiceContext struct.
//
// TODO: This implementation is prone to error if the same service sends
// MsgTick more than once, fix?
func (g *Game) waitOnServiceStart(input chan Msg) {
    // We can discard the ok value, because svc is always a struct
    val, _ := (reflect.NewValue(g.svc)).(*reflect.StructValue)
    svc_num := val.NumField() - 1 // Don't count Game
    num_left := svc_num

    for {
        msg := <-input
        switch m := msg.(type) {
        case MsgTick:
            num_left--
            if num_left == 0 {
                goto done
            }
        default:
            panic("Received message other than MsgTick!")
        }
    }
done:
}

// Starting to get somewhat game specific?

// Creates a new player entity for a requesting client
func (g *Game) spawnPlayer(msg MsgSpawnPlayer) {
    if g.PlayerFunc == nil {
        return // Skip spawning a player
    }
    p := g.PlayerFunc(g.GetUid())
    ch := make(chan Msg)
    g.ents[ch] = p
    go p.Run(ch)
    go func(uid UniqueId) {
        msg.Reply <- &EntityDesc{ch, p.Uid(), p.Id(), p.Name()}
    }(p.Uid())
}