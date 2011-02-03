// Copyright 2010, 2011 The ghack Authors. All rights reserved.
// Use of this source code is governed by the GNU General Public License
// version 3 (or any later version). See the file COPYING for details.

package core

import (
    "cmpId/cmpId"
)

// Universal message interface for components and services.
// Does not currently do anything special, but is reserved for any possible
// future use and thus should be specified where any other actual message
// types are expected.
type Msg interface{}

// Message to update
type MsgTick struct{}

// Message requesting a certain state to be returned
// Contains a channel where the reply should be sent
type MsgGetState struct {
    StateId    cmpId.StateId
    StateReply chan State
}

// Message to add an action that contains the action to be added
type MsgAddAction struct {
    Action Action
}
