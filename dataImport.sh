#!/bin/bash

# This script is meant to insert values into the db.
# This script requires the bash shell to be installed in /bin/bash which is the default on linux and mac but will not work on windows
# the cd and echo buildin commands are also required.

echo "Inserting Values"
sqlite3 data/data.db <<EOS
.mode csv

.import courses.csv Course
Update Course set DeletedAt=NULL where DeletedAt="NULL";

.import settings.csv Setting

EOS


