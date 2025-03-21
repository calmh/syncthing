-- Copyright (C) 2025 The Syncthing Authors.
--
-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this file,
-- You can obtain one at https://mozilla.org/MPL/2.0/.

--- Backing for the MtimeFS
CREATE TABLE IF NOT EXISTS mtimes (
    folder_idx INTEGER NOT NULL,
    name TEXT NOT NULL,
    ondisk INTEGER NOT NULL, -- unix nanos
    virtual INTEGER NOT NULL, -- unix nanos
    PRIMARY KEY(folder_idx, name),
    FOREIGN KEY(folder_idx) REFERENCES folders(idx) ON DELETE CASCADE
) STRICT, WITHOUT ROWID
;
