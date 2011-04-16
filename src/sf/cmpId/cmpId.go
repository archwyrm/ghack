// Copyright 2010, 2011 The ghack Authors. All rights reserved.
// Use of this source code is governed by the GNU General Public License
// version 3 (or any later version). See the file COPYING for details.

// Namespace for component indentification in order to ensure uniqueness
// in bookkeeping.
package cmpId

import core "core/cmpId"

// States
const (
    EntityList = iota + core.STATE_END
    Position
    Asset
    Health
    MaxHealth
    KillCount
)

// Actions
const (
    AddEntity = iota + core.ACTION_END
    Move
)

// Entities
const (
    Player = iota + core.ENTITY_END
    Spider
)
