// Copyright 2010, 2011 The ghack Authors. All rights reserved.
// Use of this source code is governed by the GNU General Public License
// version 3 (or any later version). See the file COPYING for details.

// Implements the publish/subscribe messaging model.
//
// This model allows asynchronus communications where senders and receivers
// are not explicit. Rather, users subscribe to certain named topics and when
// a message is published to a topic, all the subscribers receive this message.
package pubsub

import (
    .   "core"
)

type ChanType chan Msg

// Message to signal publishing of the passed data
type PublishMsg struct {
    Topic string
    Data  interface{}
}

// Message to setup a subscription to a given topic
type SubscribeMsg struct {
    Topic     string
    ReplyChan ChanType
}

// Message to remove a subscription to a given topic
// ReplyChan is to identify the subscriber
type UnsubscribeMsg struct {
    Topic     string
    ReplyChan ChanType
}

// Publish/Subscribe struct
type PubSub struct {
    *HandlerQueue
    svc           ServiceContext
    subscriptions map[string][]ChanType
    input         chan Msg
}

// Creates a new PubSub and returns a pointer to it
func NewPubSub(svc ServiceContext) *PubSub {
    hq := NewHandlerQueue()
    return &PubSub{hq, svc, make(map[string][]ChanType), nil}
}

func (ps *PubSub) Chan() chan Msg { return ps.input }

// Starts a loop to receive and handle messages from the passed channel
func (ps *PubSub) Run(input chan Msg) {
    ps.input = input
    Send(ps, ps.svc.Game, MsgTick{input}) // Service is ready

    for {
        ps.handle(ps.GetMsg(input))
    }
}

func (ps *PubSub) handle(msg Msg) {
    switch m := msg.(type) {
    case PublishMsg:
        ps.publish(m)
    case SubscribeMsg:
        ps.subscribe(m)
    case UnsubscribeMsg:
        ps.unsubscribe(m)
    }
}

// Sends a message to subscribers asynchronusly if the receiving channel blocks
func (ps *PubSub) publish(msg PublishMsg) {
    for _, sub := range ps.subscriptions[msg.Topic] {
        Send(ps, sub, msg.Data)
    }
}

// Adds a subscription to the appropriate topic
func (ps *PubSub) subscribe(msg SubscribeMsg) {
    subscribers := ps.subscriptions[msg.Topic]
    ps.subscriptions[msg.Topic] = append(subscribers, msg.ReplyChan)
}

// Removes a subscription from the given topic
func (ps *PubSub) unsubscribe(msg UnsubscribeMsg) {
    subs := ps.subscriptions[msg.Topic]
    var rm_i int
    for i, s := range subs {
        if msg.ReplyChan == s {
            rm_i = i
            break // TODO: Remove multiple or disallow multiple subscription?
        }
    }

    // Slice around rm_i
    subs = append(subs[:rm_i], subs[rm_i+1:]...)
    ps.subscriptions[msg.Topic] = subs
}
