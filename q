[33mcommit f5f6bdd8bbd8250ac83b30da1cce2be8458b1192[m[33m ([m[1;36mHEAD -> [m[1;32mfeature/database-add-tests[m[33m)[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Mar 9 16:49:54 2024 +0100

    Testing: Moved tests out of subfolders into tests/
    
    This was done, because previously the automatic testing tool could not find the tests in the subfolders.

[33mcommit b7cd500bb8ccb8b2aee0b6778b25c8e063759a26[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Mar 9 16:48:55 2024 +0100

    Testing: Fixed `...Equal` functions and added settings tests.

[33mcommit 4406d0692b38c6bc2bdc2e9f18e8dc7f60c35dc3[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Mar 9 16:35:49 2024 +0100

    Database: New tests for general database tasks.

[33mcommit 704cf3e67125578492b2dcc2adb923bca2485fc2[m
Merge: 824d020 9f202b9
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Mar 9 15:55:01 2024 +0100

    Merge "feature/frontend/hmi_backend_integration", because of updated working `CourseStudents` function.

[33mcommit 9f202b9773a48769f84746e5f66412dd5dab1621[m[33m ([m[1;31morigin/feature/frontend/hmi_backend_integration[m[33m)[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Thu Mar 7 23:37:57 2024 +0100

    Working Overview, Student, Course Page with Database Implementation

[33mcommit d0607ef7e194c593df3eac0cbdaeb381ecccc9ae[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Wed Mar 6 20:02:45 2024 +0100

    Compilable Version

[33mcommit 74e5a7a32b524dc8a65020d214af8b0716db0691[m[33m ([m[1;31morigin/feature/backend-e-mail-processing-thread[m[33m, [m[1;32mfeature/backend-e-mail-processing-thread[m[33m)[m
Merge: 5e5e3d5 54624ec
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Mar 6 15:08:32 2024 +0100

    Merge remote-tracking branch 'origin/dev' into feature/backend-e-mail-processing-thread

[33mcommit 54624ec4819ee8e8fbb4569dfd07323f9f11074a[m
Merge: 4ac0b2c 8104bc9
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Mar 6 15:29:27 2024 +0100

    Merge pull request #85 from DHBW-SE-2023/bugfix/build-on-macos
    
    Building on macos with tesseract and leptonica

[33mcommit 5e5e3d575f56970c20fe66e2f60e1a2e602180be[m[33m ([m[1;31morigin/feature/hmi_backend_integration[m[33m)[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Mar 6 15:01:12 2024 +0100

    ReceiveNewTable in MVVM Frontend part now gets called when NotifyNewList is called

[33mcommit 687588126caecc16ce72445aea7d7b3ed77361ac[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Mar 6 14:52:59 2024 +0100

    Frontend can call MVVM now

[33mcommit eb0bd3efa910f296266cc747fc7c6b8671d7181f[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Mar 5 11:08:41 2024 +0100

    Demon: After every successfully parsed list the frontend gets notified.

[33mcommit 0ef876b0762f181f85c6f1af2c65ce1cf5afdec0[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Mar 5 10:56:49 2024 +0100

    Added notification function to MVVM for the Frontend to hook into.

[33mcommit cd45846a9a1496c435d7f5f1ec9d3a7855800dc3[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Mar 5 10:30:19 2024 +0100

    Use of the actual mail functions in the demon.

[33mcommit 2ec4546237b0d96b7354db9520e1b56c8eccea79[m
Merge: dad55db 4ac0b2c
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Mar 5 09:58:52 2024 +0100

    Merge branch 'dev' into feature/backend-e-mail-processing-thread

[33mcommit dad55db5f0bd6b4c1d84ed43e9ea4ed5d11348be[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Mon Mar 4 15:05:42 2024 +0100

    Added a background demon to fetch new emails and automatically process them

[33mcommit 824d020d4008eca2a67aee62171e16c3de771382[m[33m ([m[1;31morigin/feature/database-add-tests[m[33m)[m
Merge: 0a7a253 4ac0b2c
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Mon Mar 4 12:41:21 2024 +0100

    Merge pull request #80 from DHBW-SE-2023/dev
    
    Make main great again

[33mcommit 8104bc94b70c2832642a4bfef05bcf23f385f7d0[m[33m ([m[1;31morigin/bugfix/build-on-macos[m[33m)[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Mon Mar 4 10:15:24 2024 +0100

    Now with the libraries tesseract and leptonica installed, the build system works again.

[33mcommit 4ac0b2c1fec29719cecdb2fbadef474544d9cb40[m[33m ([m[1;32mdev[m[33m)[m
Merge: 95ad226 2172435
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Sat Mar 2 16:23:26 2024 +0100

    Merge pull request #79 from DHBW-SE-2023/docker_image
    
    Add docker-image build to dev

[33mcommit 2172435697b5f59d279b2a03e68edfb51fa8a614[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Sat Mar 2 16:22:56 2024 +0100

    Added tesseract deps

[33mcommit 95ad226e88c5945bf6da77ce63789f26b6545c5d[m
Merge: a25465c 894055f
Author: Vinzent <114230986+V1nzent@users.noreply.github.com>
Date:   Sat Mar 2 16:05:02 2024 +0100

    Merge pull request #76 from DHBW-SE-2023/Mail_attachment_extraction
    
    Mail attachment extraction

[33mcommit 894055f3e881ed13a5cbfbcb4a680ebc3d5ca1b6[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sat Mar 2 15:59:04 2024 +0100

    changed test function

[33mcommit a25465cc6bd356a4432f351eb9a6bc88d6165bad[m
Merge: 99184cd 7d793e7
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Mar 1 14:30:32 2024 +0100

    Merge pull request #77 from DHBW-SE-2023/feature/mvvm-implement-mvvm-for-main-window
    
    Deletes all your code, refuses to elaborate further, leaves

[33mcommit 7d793e7a8515855369ce724d390159a062d82df7[m[33m ([m[1;31morigin/feature/mvvm-implement-mvvm-for-main-window[m[33m, [m[1;32mfeature/mvvm-implement-mvvm-for-main-window[m[33m)[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Mar 1 12:47:22 2024 +0100

    More database functions and tests

[33mcommit c8cced388d07274f74ef665784cb5e6bb3cc4ef8[m[33m ([m[1;31morigin/Mail_attachment_extraction[m[33m)[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Mar 1 10:38:58 2024 +0100

    functions so it works

[33mcommit fbe7e11310f9ad1212f6d477cea53539654cd659[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Mar 1 10:38:11 2024 +0100

    commented

[33mcommit 4b8f1eb74237a43e1c5ebfea4719940529252676[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Feb 28 00:38:36 2024 +0100

    Testing: Moved Imgproc tests into a more sensible folder

[33mcommit 37404eb358869689624b609a8e6a04341188100b[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Feb 24 10:44:03 2024 +0100

    Added the functionality as described by Eric's overview for the required functionality.

[33mcommit 62171d0350e082f5721fa36aa476f0a4ada16f20[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 16:29:31 2024 +0100

    Deletes all your code, refuses to elaborate further, leaves
    
    As I was interested how it would turn out, I rewrote the database part using GORM as the
    ORM of choice. As it turned out This cut down tremendously on the size and therefore
    cognitive overload of the database package (also because I did no error handling).
    I rewrote some functions which could now be implemented in a much shorter fashion.
    Those functions were also made available to the MVVM.

[33mcommit 5c64c207707ea25badfefb19094ec88cc559e8a9[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 23 15:38:44 2024 +0100

    added error if no boundary found

[33mcommit f106aa503a761489ff04af35f50f506488cb67fe[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 23 14:44:51 2024 +0100

    added comments

[33mcommit 07210a797dfa9e6bc9ab626d6a65b529d987e94f[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 23 14:44:03 2024 +0100

    removed unnecassary method

[33mcommit c3488ea605c7ec59dc804ac6bebecdd9d9cc3961[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 23 13:46:11 2024 +0100

    small fix

[33mcommit 78c55a41f09ab901be5f9f6fc2eb49a48c5b3941[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 23 13:44:58 2024 +0100

    added nonTLS

[33mcommit 99184cd12654380e68c1cea590464f81a1b6c158[m
Merge: 5dd49be c790916
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 12:19:16 2024 +0100

    Merge pull request #74 from DHBW-SE-2023/feature/backend-make-gocv-code-available-to-the-mvvm
    
    MVVM: make imgproc code available to the mvvm

[33mcommit a3b3ac1b9a2de7e56dfa5b06fddcf6c327037d6e[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 23 12:13:53 2024 +0100

    added func for testing

[33mcommit 5dd49be16d06ac4e443221d191c51b100214a9e3[m
Merge: 2630a9f 3beecd7
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 12:07:14 2024 +0100

    Merge pull request #72 from DHBW-SE-2023/feature/image-processing-bounding-box-gathering-of-each-name-signature-cell
    
    Image Processing: Bounding box of rows and refactor of GoCV code

[33mcommit c7909164c91b60e51ba3874191eb8425ca3b628f[m[33m ([m[1;31morigin/feature/backend-make-gocv-code-available-to-the-mvvm[m[33m, [m[1;32mfeature/backend-make-gocv-code-available-to-the-mvvm[m[33m)[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 11:48:47 2024 +0100

    MVVM ValidateTable calls backend and on success notifies frontend.

[33mcommit 7f6bd720eaea44ce11f206978814593755c321dd[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 10:59:22 2024 +0100

    Added `ValidateTable` functionality for Imgproc MVVM backend

[33mcommit 9becc3fab80ea4279662b3dc0f075c1c36f32c24[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 10:39:54 2024 +0100

    Moved CV/Imgproc code into the backend directory
    
    I moved the CV code (now called Imgproc, because that's better fitting) into the
    internal/backend/imgproc folder, as a start to integrate it with the MVVM.
    For this I also removed the old test code in internal/backend/opencv.

[33mcommit 3beecd7fe1206cdf38481b6cb91bcf45e07ba259[m[33m ([m[1;31morigin/feature/image-processing-bounding-box-gathering-of-each-name-signature-cell[m[33m)[m
Merge: 5eac0a5 2630a9f
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 10:15:07 2024 +0100

    Merge branch 'dev' of github.com:DHBW-SE-2023/YAAC into dev

[33mcommit 5eac0a50321172ec595043db3648bf0cc86a975e[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Feb 23 10:09:50 2024 +0100

    CV code refactoring: fixed merge, full row bounding rect
    
    I mainly refactored the code in various places to return a more generic type `TableRow`
    instead of, for example, `NameROI`. This required also a rewrite of the merge function
    where I realised some errors of it and at the end managed to refactor it into a simpler,
    but also more readable form.

[33mcommit 43d77198adbca17af179331f82bb25dbbd1dd05f[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 23 10:06:34 2024 +0100

    removed unnecessary stuff

[33mcommit 2630a9f87e5ca496b9624d737de858770119f11c[m
Merge: 2471e2d 613b41f
Author: Jonas Zagst <87520165+JonasZagst@users.noreply.github.com>
Date:   Fri Feb 23 08:22:25 2024 +0100

    Merge pull request #71 from DHBW-SE-2023/databasesetup
    
    added requested database functions and made path variable

[33mcommit 613b41f11c99fa8058dbfeafb0b2962f4c312eea[m
Author: JonasZagst <jonas@zagst.net>
Date:   Wed Feb 21 11:21:44 2024 +0100

    removed exclusion of /test/testdata

[33mcommit fa3c45694adeb7a8433b3fe78d16f55eec1f616c[m
Author: JonasZagst <jonas@zagst.net>
Date:   Wed Feb 21 10:52:12 2024 +0100

    updated GetAllAttendanceWithStudentName to also return the day

[33mcommit 966972c97b71f810ea530a8f90d6376074e186e2[m[33m ([m[1;31morigin/docker_image[m[33m)[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Wed Feb 21 10:09:53 2024 +0100

    Install llvmpipe

[33mcommit 310bbd3db5aaebf8efc65f84280b824d3c8e3ceb[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Mon Feb 19 08:24:36 2024 +0100

    added getCourse

[33mcommit 5e6e0720be13098e835dc74a8493993092fed45a[m
Author: JonasZagst <jonas@zagst.net>
Date:   Sun Feb 18 23:30:22 2024 +0100

    updated insert attendance

[33mcommit 99864eb411f534b624103b30065e18c790c27d87[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sun Feb 18 16:33:16 2024 +0100

    comments

[33mcommit 5ce6ba547f2560d349d666dfb1c19d6b810ede5b[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sun Feb 18 13:31:16 2024 +0100

    comments

[33mcommit 22c5510a0f4412f0b3533a5cb5d473985fb72aea[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Sun Feb 18 11:47:20 2024 +0100

    Fixed trigger name

[33mcommit 08546496dff4e032b6d8910df24c6ccbcf9cc13b[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Sun Feb 18 11:46:30 2024 +0100

    Removed old docker action

[33mcommit 696986b34511c2a9632e8888536ff2a0b912453f[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Sun Feb 18 11:40:16 2024 +0100

    Create docker-publish.yml

[33mcommit d79d0f1308a81157bb59c7d27bb2f012a7f0df2f[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Sun Feb 18 11:30:53 2024 +0100

    Update docker.yml

[33mcommit d12bc905f93dbd9aaf53e429e95b50c82f09d219[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Sun Feb 18 11:29:59 2024 +0100

    Update docker.yml

[33mcommit e5bda743b0aec7edbd55b9843063fb463ba439a8[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sat Feb 17 21:52:07 2024 +0100

    added Date checker

[33mcommit fc3f05e2855c2a6fc8350d35a2d18999c6ab0ebc[m
Author: JonasZagst <jonas@zagst.net>
Date:   Fri Feb 9 18:07:05 2024 +0100

    added requested database functions and made path variable

[33mcommit 8cf8f060a5ceabe17c3eaffd4a98263ae6963729[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Fri Feb 9 13:53:13 2024 +0100

    Escaped $ in makefile

[33mcommit 38fc9a22f23c4efffe1a9d64a5ead0077bab96cf[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Fri Feb 9 13:53:00 2024 +0100

    Hardcoded base domain in dockerfile

[33mcommit db3c383f7c7df17805169e8ebddbe662a01a7fc4[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 9 12:26:05 2024 +0100

    added getCourse and getDateTime Template

[33mcommit a2c5ca9f4daef4f27e0bd5b0cf90b529bae70629[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Fri Feb 9 12:11:30 2024 +0100

    Added deps

[33mcommit 443c57a88414ebc056ba3d7aa29f95843096fce5[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Thu Feb 8 13:29:07 2024 +0100

    Added full build pipeline

[33mcommit 77201f5bb93147cf0ddebce8bcbfbc911200bda4[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Thu Feb 8 11:58:37 2024 +0100

    Added basic docker toolchain

[33mcommit 2471e2d972ed0d88d68d8b53acdfe9648d60f402[m
Merge: f2d72bd 968a4e1
Author: danielSiegt08 <124031951+danielSiegt08@users.noreply.github.com>
Date:   Tue Feb 6 09:15:00 2024 +0100

    Merge pull request #67 from DHBW-SE-2023/feature/hmi-fyne-style
    
    Feature: Implement HMI fyne style and layout

[33mcommit 968a4e181861a6028e5680c8d77a8b120dc78030[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Tue Feb 6 08:27:56 2024 +0100

    Added settings nav functionalities, as some minor layout changes

[33mcommit f2d72bdc8f1c6dc362f90131ae7aa4b797ebba04[m
Merge: 4625b7e 0215fb4
Author: Jonas Zagst <87520165+JonasZagst@users.noreply.github.com>
Date:   Tue Feb 6 08:20:43 2024 +0100

    Merge pull request #69 from DHBW-SE-2023/databasesetup
    
    Implemented basic database functions and the data infrastructure

[33mcommit 0215fb4e6289db866455359f5cceccce5d0fb677[m
Author: JonasZagst <jonas@zagst.net>
Date:   Mon Feb 5 23:22:37 2024 +0100

    added several database functions for frontend

[33mcommit 967f5a9a693fe6f7cb2b03146f055aa48752db1b[m
Author: JonasZagst <jonas@zagst.net>
Date:   Mon Feb 5 13:57:06 2024 +0100

    implemented variable path to database

[33mcommit 9c4868f7147b64bc3433918baefd767ead035a55[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Mon Feb 5 10:17:04 2024 +0100

    Removed title & intro from page

[33mcommit 682fda357aac57d160ab74d9cf590bebc84ad6dd[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Mon Feb 5 09:21:56 2024 +0100

    changed user crudentials

[33mcommit a3c4d57815d4b9d91232af23649be72a9d1b7031[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Feb 5 08:55:16 2024 +0100

    Removed mobile references

[33mcommit 6d22c0fe7493d2266813a747f427e61d000af362[m
Merge: a680ed2 4625b7e
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Feb 5 08:49:31 2024 +0100

    Merge branch 'dev' into feature/hmi-fyne-style

[33mcommit a680ed2fc25ed3ea2cdbb930749279f469cb3cea[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Feb 5 08:46:48 2024 +0100

    Fixed transparent dropdowns

[33mcommit 698fba6f70c79c29d9ea264e9db27dfc4940b776[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Feb 5 08:43:51 2024 +0100

    Re-Implemented sidebar images

[33mcommit 6bac93fc719f81ea4b9aa386e6693b77ad777a6f[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Feb 5 08:37:02 2024 +0100

    Reverted custom window style changes

[33mcommit 8c6fa2fcdef1f8a0b3c74786c799838a55cf78b9[m
Merge: 691e279 4a62d48
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Feb 5 08:31:24 2024 +0100

    Merge branch 'feature/hmi-fyne-layout' into feature/hmi-fyne-style

[33mcommit 4d8241bfabf5e5bb0d6c5dc1374937886d75e309[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Mon Feb 5 08:11:33 2024 +0100

    small fixes

[33mcommit 8f00603fbbd38aa2e0c21099b7559760d46aaf5a[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sun Feb 4 12:47:37 2024 +0100

    comments, integrated base64decode

[33mcommit 5c310352e4ab9bab8edb1c16b6fb4dddadd369eb[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sat Feb 3 19:28:41 2024 +0100

    removed unnecessary stuff, added methods

[33mcommit 90c3d2d4816cd6ddbaf78a7fd59eae90964381b4[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 2 19:33:50 2024 +0100

    changed to structs

[33mcommit 09c4987cabeb1e030edf01ec93122eb0d637caea[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Feb 2 19:33:14 2024 +0100

    update New() backend mail

[33mcommit d175d47712a2df3dae5cead79e7ac1177dbfa180[m
Merge: 2701fc6 4625b7e
Author: JonasZagst <jonas@zagst.net>
Date:   Fri Feb 2 16:59:36 2024 +0100

    Merge remote-tracking branch 'origin/dev' into databasesetup
    
    # Conflicts:
    #       cmd/yaac/main.go

[33mcommit 691e279772cd29e232777baab62e4cc9f04549f6[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Fri Feb 2 14:34:55 2024 +0100

    Added click event to button

[33mcommit a6949521f0638e81a3f2a76500cb252f5d878a48[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Fri Feb 2 11:56:45 2024 +0100

    It's complicated

[33mcommit 4a62d48064b48345a793bd15d5ffa05f1eda1d46[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Thu Feb 1 16:53:16 2024 +0100

    Small Layout Changes, General Settings Layout done, Functionalities missing

[33mcommit 1eeac0d3dd1b38c229f8f7fa78412f2db8097cf4[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Wed Jan 31 17:22:16 2024 +0100

    Implemented First Layout for Student View + Pop Up Window

[33mcommit 36e0930af18a0620e4c0998a10ba880eb6dd5523[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Wed Jan 31 16:17:40 2024 +0100

    Added Custom Theme

[33mcommit 2701fc6a2babd7166129202e3f71994eaef24f34[m
Author: JonasZagst <jonas@zagst.net>
Date:   Wed Jan 31 00:09:15 2024 +0100

    fixed database.go

[33mcommit da15f536b8e3ffef6ef8303545214fa03149ca5e[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Tue Jan 30 17:19:04 2024 +0100

    First Version of Overviwe Page

[33mcommit 78d68ad027a3597cad97551a05011bc6b2196830[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Jan 30 17:01:11 2024 +0100

    Refactoring code into yaacs codestyle

[33mcommit 4625b7eb423189d272411169d3460d7fa65c394f[m
Merge: 560fb90 1077560
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Jan 30 15:18:29 2024 +0100

    Merge pull request #58 from DHBW-SE-2023/feature/image-processing-signature-verification
    
    Feature/image processing signature validation

[33mcommit 1077560a851590c361832a1ca3bbfdae3598a824[m[33m ([m[1;31morigin/feature/image-processing-signature-verification[m[33m, [m[1;32mfeature/image-processing-signature-verification[m[33m)[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Jan 30 15:07:43 2024 +0100

    Pulled merging logic out of cv.ValidSignature into a helper function

[33mcommit 8a3bd37062de467ee1af150eac5724dad9a473d0[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Jan 30 15:02:12 2024 +0100

    Added case in cv.ValidSignature for if r2 is left of r1

[33mcommit bd74b048881042708703ca56f734729734dd4457[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Jan 30 14:39:20 2024 +0100

    Added functionality in cv.ValidSignature to merge bounding rectangles
    
    Sometimes it happens, that the contour detection splits a signature apart into two
    bounding rectangles. This could then result in the signature field getting marked
    as invalid. As an example for where this could happen, take a look at the signature field
    in the row designated "11".
    To mitigate such issues, I made it, such that after the bounding rectangles were gathered based
    on the contours of the signature, those bounding rects were merged together into one rectangle,
    which is the bounding box of both former rectangles.
    After that the previous algorithm is applied on the merged bounding rectangles.

[33mcommit 570e7ff91dfc8a7974566f6ade7f2e61d2030d0b[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sat Jan 27 13:55:19 2024 +0100

    added comments

[33mcommit c4f2cdfffa6fcc7b250133cd3880517c7cefc9e5[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Jan 26 20:40:44 2024 +0100

    writes images to home directory

[33mcommit 83692580bb38e1ccbb9880f69c3f0a5d899e7449[m
Author: danielSiegertUAM <daniel.siegert@airbus-uam.com>
Date:   Fri Jan 26 19:26:55 2024 +0100

    First Implementation of Sidebar

[33mcommit f000082d26583a723e57a4d7cbc384ef087a5070[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 26 10:53:33 2024 +0100

    Added tests for cv.ReviewTable
    
    This test is to ensure that cv.ReviewTable performs all the correct logic.
    This is added instead of cv.ValidSignature, because cv.ValidSignature can only be called
    in a certain environment, which is set up by cv.ReviewTable. There would be no sense in
    recreating that environment in a test function, if it is already created by cv.ReviewTable.
    Therefore we simply use cv.ReviewTable.

[33mcommit 940028d66d3325fd5e155bf250094a755eec7f0f[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 26 10:12:14 2024 +0100

    Image processing: signature validation
    
    Signature validation was achieved in a very simple way, that should be
    subject to change in the future. We do this validation by simply asserting
    that each signature field only has one signature in it. This is done by
    detecting the contours and getting the bounding boxes of them. We select
    the interessting bounding boxes and check that there is only one of them
    per signature field.
    We should replace signature validation with signature verification where
    we also assert that the signature belongs to the correct student.

[33mcommit c0b16bb6861782a3fe94ea196c04feeb15331c3e[m
Author: JonasZagst <jonas@zagst.net>
Date:   Fri Jan 26 00:02:37 2024 +0100

    basic implementation of sqlite database

[33mcommit 560fb90a149a8c0daeaf281e6d3beedb0987c464[m
Merge: cc43444 930348e
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 19:38:05 2024 +0100

    Merge pull request #49 from DHBW-SE-2023/image-processing-get-the-student-names-from-the-table
    
    Image processing get the student names from the table

[33mcommit 930348ec7093e2d22c71de1db744a625d97b7681[m[33m ([m[1;31morigin/image-processing-get-the-student-names-from-the-table[m[33m, [m[1;32mfeature/image-processing-cell-extraction-out-of-table[m[33m)[m
Merge: 3a672ba dbe822d
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 17:49:33 2024 +0100

    Merge branch 'image-processing-get-the-student-names-from-the-table' of github.com:DHBW-SE-2023/YAAC into feature/image-processing-cell-extraction-out-of-table

[33mcommit 3a672bae3d6c1189bdd97247a4e262047bf7e5dc[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 17:49:29 2024 +0100

    Fixed test which broke as a result of moving CvtColor out of NewTable.

[33mcommit dbe822d72d581f102bb40363d5136c2ac8f7644f[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 17:43:21 2024 +0100

    Update test pipeline: Install OCR Tesseract

[33mcommit e38068617bb837514102d7f5e7b042e6d72b484a[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 17:34:46 2024 +0100

    Improved accuracy of text recognition by inverting image and not being retarded

[33mcommit 179d25edcd669870f3f6400e0f11b798bf19bef8[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 16:12:03 2024 +0100

    Removed initial color conversion from NewTable
    
    The color conversion from BGR to grayscale is done by the caller.

[33mcommit 8353e948c7d07aa43e4b0392fee833e65e2625aa[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 16:11:26 2024 +0100

    Added getters for NameROI.name and NameROI.roi

[33mcommit bfbff66b247f220dd481cdcd1c3fcf6d74dbf52c[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 15:36:27 2024 +0100

    Refactored cv.ReviewTable, extracted character recognition
    
    Through refactionring `ReviewTable` I extracted the name recognition in its own function.
    It expects a preprocessed image on which it uses OCRTesseract.

[33mcommit d390833b5452d7fabb022920e9c71fcf3c8643cd[m
Merge: 064f7df cc43444
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 14:33:16 2024 +0100

    Merge branch 'image-processing-get-the-student-names-from-the-table' of github.com:DHBW-SE-2023/YAAC into feature/image-processing-cell-extraction-out-of-table

[33mcommit 064f7df317cebde6e08ba02f08ff4c2a6bb594d2[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Wed Jan 24 14:28:29 2024 +0100

    cv.ReviewTable extracts the students' names from the name column.
    
    This is achievied by first preprocessing the image. After that a NewTable applied
    on the image, which returns a table based on the image. After that we iterate
    through all rows and skip any deformed rows. NewTable already performs some checks
    to insure that this doesn't happen, this is just to be safe.
    Afterwards be extract the ROI of the name region and convert it to PNG. This is done
    because OCRTesseract only accepts images of certain formats like PNG or JPEG.
    We use PNG because it is a lossless format and therefore no additional recognition
    errors are created through the lossy compression.

[33mcommit e8d21557465c11a9d07f8a5da4424765184d8244[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Sun Jan 21 18:04:01 2024 +0100

    added mail service

[33mcommit cc4344464df42ea9261e67e3d7ab21bd3d839268[m
Merge: bfbdb8b d063343
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 21:50:48 2024 +0100

    Merge pull request #44 from DHBW-SE-2023/feature/image-processing-cell-extraction-out-of-table
    
    Image Processing: cell extraction out of table

[33mcommit d063343c4e69e3dcdf5df20ef75270cdd8a2b99a[m[33m ([m[1;31morigin/feature/image-processing-cell-extraction-out-of-table[m[33m)[m
Merge: 584903d 99a0851
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 19:17:21 2024 +0100

    Merge branch 'feature/image-processing-cell-extraction-out-of-table' of github.com:DHBW-SE-2023/YAAC into feature/image-processing-cell-extraction-out-of-table

[33mcommit 99a0851545a4c627d9ea3d927682d14fd7a80bb7[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 19:06:01 2024 +0100

    Set access token ... again
    
    Why is this not working? Must have copy-pasted it wrong the first time, but I try again.

[33mcommit 71719beff63dbd0e502d6eac7149aa6be81d6b67[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 18:41:07 2024 +0100

    Add an PAT to the pipeline to access the private testing-data repository

[33mcommit 40b75eb9c1ed057865a922d9f06eec5cb035d10e[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 18:36:32 2024 +0100

    Checking out the code in the pipeline now also fetches submodules
    
    Because we are using Git Submodules to handle private data we also need to change the pipeline, s.t. when checking out the code it also clones the submodules.

[33mcommit 584903d48071b7367240d9c02e6bd6c89f8f6bfd[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 18:25:15 2024 +0100

    Change path for image in TestTableColumnCount from test/testdata/... back to testdata/...

[33mcommit 125247297e1f4c900d1b2a63692639fd1b75ac1d[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 18:21:41 2024 +0100

    Added new check if the image in TestTableColumnCount is empty, due to some problems with the test pipelnie in GitHub.

[33mcommit 7958c1a0b8a4e90bdc6a2b20ad96b37752367ccd[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 18:15:41 2024 +0100

    Changed test to open use test/testdata/ instead of data/

[33mcommit 263db05bf5a66957517a327f114250fbeb23590e[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Sat Jan 20 18:14:19 2024 +0100

    Moved data/ to test/testdata to comply with Go conventions regarding data used by tests.

[33mcommit bfbdb8bdfdbd4ab94d4f0bf4e75aaaa67e7c7cae[m
Merge: 86b9949 0a3c426
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 22:55:36 2024 +0100

    Merge pull request #42 from DHBW-SE-2023/feature/testing-private-repository-for-sensitive-data
    
    Testing: private repository for sensitive data

[33mcommit 6488a640b289078e1c67a6a86fbfafd26a94e9bf[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 22:41:52 2024 +0100

    Testing: Assert that an attendance sheet has exactly 3 columns.

[33mcommit 8f9f125e2c899bd363d4354e5a30fc6d5a348c05[m
Merge: fc19088 bfbdb8b
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 22:30:29 2024 +0100

    Merge branch 'dev' into feature/image-processing-cell-extraction-out-of-table

[33mcommit 0a3c426e52e76c238a0bc2ae8807b14c1d99a781[m[33m ([m[1;31morigin/feature/testing-private-repository-for-sensitive-data[m[33m, [m[1;32mfeature/testing-private-repository-for-sensitive-data[m[33m)[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 22:09:12 2024 +0100

    Added image of attendance list to data/

[33mcommit a85721657ed21f2145f8a4776d655f7d2d43fab3[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 22:05:18 2024 +0100

    Added Git Submodule for using sensitive private data in tests and development

[33mcommit fc19088c5ab69a112b17c8dbc07246a883a931c1[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 22:02:25 2024 +0100

    Removed redundant check if the width or height of a cell is zero.

[33mcommit 67309b02ef2f39b36732fa4cf8b5d73d3dd60e04[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 21:23:35 2024 +0100

    Additional checks for parsing a table
    
    An additional check when parsing a table was added, that checks
    whether the column that should be added has a minimum height and a minimum width.
    It turns out that the following to functions to compute the minWidth and minHeigh turn out
    to work pretty well:
    minWidth = int(minWidth * 0.04)
    minHeight = int(maxHeight * 0.02)

[33mcommit ee741f50e66de53fd5a0dd80bc27bbb48bfd5205[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Jan 19 21:02:36 2024 +0100

    writing image data to file

[33mcommit f2b221fc79c2d3d58898642750520ad557ff6d7c[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 20:58:36 2024 +0100

    Added image preprocessing to `NewTable`
    
    Added image preprocessing to `NewTable` to match the preprocessing done in `FindTable`.
    This is done because I suspect that in the future we change the input image format of `NewTable`,
    so I did it now. So `FindTable` now returns an RGBA image.
    Because of UI reasons we probably want to return the original image with the perspective transform applied to it,
    instead of the binary version of the original image.
    Because we then take in the normal image in `NewTable` we need to apply the same image manipulations to it as in `FindTable`.

[33mcommit d8056989a5c2a625c2bdc0e754be8adc117e7bab[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 20:36:20 2024 +0100

    Find the cells of the table put them into a row-major matrix

[33mcommit 10185822b2e73a4a90c9aaea72a2be75a76143c3[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 20:31:43 2024 +0100

    Extract the vertical and horizontal lines from the table image
    
    Extract the vertical and horizontal lines from the table image and recombine them into a new image consisting only of the vertical and horizontal lines.

[33mcommit 668d9433a04d13f2624dd2befc8fbfd43d5d614a[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Jan 19 17:02:56 2024 +0100

    added comments

[33mcommit 86b99497c204b085767479295393e24b23f7b1ee[m
Merge: eae944b 95a4b84
Author: Erix0815 <80325429+Erix0815@users.noreply.github.com>
Date:   Fri Jan 19 16:58:48 2024 +0100

    Merge pull request #40 from DHBW-SE-2023/feature/find-table-warp-perspective
    
    Feature/find table warp perspective

[33mcommit 95a4b844dcf1c88aad9db87f15f080306cc52e74[m[33m ([m[1;31morigin/feature/find-table-warp-perspective[m[33m, [m[1;32mfeature/imgproc-submodule[m[33m)[m
Merge: 9183be8 eae944b
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 16:34:07 2024 +0100

    Merge branch 'feature/find-table-warp-perspective' of github.com:DHBW-SE-2023/YAAC into feature/imgproc-submodule

[33mcommit 9183be8ef79277f34ef3849d6996383c8ed9ccc3[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 16:23:32 2024 +0100

    Return a prepared version of the image, without the Canny filter applied to it

[33mcommit 1d55cf2b8424caad602f2d73cb07484b8d4c6b18[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 16:21:37 2024 +0100

    Find the largest rectangle and warp the perspective to only show it

[33mcommit eae944b421c4d767a69b515864942a1a3ee9fad3[m
Merge: f053301 80c062e
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 16:21:36 2024 +0100

    Merge pull request #38 from DHBW-SE-2023/feature/imgproc-submodule
    
    New module to contain CV code

[33mcommit 5b5182c7ada9a7834af67c3aa17df0b243a64e56[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 16:18:40 2024 +0100

    Preprocess an image to ease further processing

[33mcommit 80c062ed133ffe734441f62e24e669b6b86bebfd[m[33m ([m[1;31morigin/feature/imgproc-submodule[m[33m)[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Fri Jan 19 15:54:44 2024 +0100

    Created a new module which contains the CV code

[33mcommit aa66d6337584e92369459b63c274bb9a3e276ddf[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Jan 19 14:59:00 2024 +0100

    added error when no image found

[33mcommit 900f8cbae0a8b029e1a2a8ef12c63360a7d8d16b[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Fri Jan 19 13:57:49 2024 +0100

    added mail attachment extraction

[33mcommit f053301fba1d4af866993d219b0be21497187b6e[m
Merge: 1013708 0a7a253
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Mon Jan 15 11:31:20 2024 +0100

    Merge pull request #34 from DHBW-SE-2023/master
    
    merge master in dev, so its up to date

[33mcommit 0a7a2533c41e1d597f368dde761714e1ca41d571[m[33m ([m[1;31morigin/feature/get-the-student-names-from-the-table[m[33m)[m
Merge: d3f54c9 1013708
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Mon Jan 15 10:54:29 2024 +0100

    Merge pull request #33 from DHBW-SE-2023/dev
    
    Make master great again

[33mcommit 7e9b87b9938d6ff8273e8a76bd90cbd8dc54d8c1[m
Author: JonasZagst <jonas@zagst.net>
Date:   Mon Jan 15 00:18:03 2024 +0100

    added database setup

[33mcommit 10137085654f86452c61edf38e4d9cf439223578[m
Merge: 4817047 ba607cb
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Nov 28 13:09:21 2023 +0100

    Merge pull request #32 from DHBW-SE-2023/Add-license-to-readme
    
    Update README.md with license

[33mcommit 481704700a38b611425aaa20b2669e1201be44bf[m
Merge: e0d2e13 4640526
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Tue Nov 28 13:08:26 2023 +0100

    Merge pull request #30 from DHBW-SE-2023/cicd_test
    
    Add pipeline for Github Actions

[33mcommit ba607cb2bdd76d0ce3003d399f3e9b999e230828[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Tue Nov 28 12:08:50 2023 +0100

    Update README.md with license

[33mcommit d3f54c9fd8d2149d4a3468ae88c623c6583072f9[m[33m ([m[1;32mmaster[m[33m)[m
Merge: 2bc9fc2 f4de830
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Nov 28 12:01:19 2023 +0100

    Merge pull request #17 from DHBW-SE-2023/licence
    
    Create LICENCE

[33mcommit e0d2e13c6d8908beb362aaa50e8ee0d86bce2d82[m
Merge: 76757e6 f4de830
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Nov 28 11:58:23 2023 +0100

    Merge pull request #31 from DHBW-SE-2023/licence
    
    Create LICENCE (in dev too)

[33mcommit f4de830da06042d3ddaad7fa35feb2f5f9eb9803[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Tue Nov 28 11:53:58 2023 +0100

    Create LICENSE GPLv3

[33mcommit 7a243dd25c60756aa034e9f0e4f0dcceab3aa38c[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Tue Nov 28 11:50:42 2023 +0100

    Delete MIT LICENCE

[33mcommit 4640526a781acba972efb008d8e25edd6fa90e7e[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:59:42 2023 +0100

    Removed debug statement

[33mcommit df7390b5d8de292e20222c4c431ea15bf5574ee8[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:58:51 2023 +0100

    I don't even know any more

[33mcommit fed4d2160b07cb4791ecf1bbc7c2a27dbbde1802[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:54:49 2023 +0100

    Trying my best to debug

[33mcommit 989d006ed64de84de60df2daf98d71cb22a4e048[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:51:50 2023 +0100

    Adapted paths to runner FS

[33mcommit 8e37f9c800b4abd97a6c7fdf1d8041e8ca963766[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:47:29 2023 +0100

    The final sudo

[33mcommit 5536f28be19305b84e06b437142f552cc0c379af[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:45:53 2023 +0100

    Added even more sudo

[33mcommit b5429ffb3188a29305c4ff741ef66a9ea6e5d4bb[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:43:21 2023 +0100

    Added sudo to apt. 'docker runs the same on any system'... as if

[33mcommit 03de61028fe6272ebacb49d3c5cc506d5d2ec931[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 28 10:10:44 2023 +0100

    Added pipeline working in act. Go pkg not caching

[33mcommit a357ee24ff9b205a49ef9468eb6fb0b2bd5980bf[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Nov 27 11:17:25 2023 +0100

    Static linking trial

[33mcommit 76757e6b049710be0d5bdf9ae863d6682c52d21e[m
Merge: 08492c3 5ff0677
Author: Jonas Zagst <87520165+JonasZagst@users.noreply.github.com>
Date:   Thu Nov 23 14:32:10 2023 +0100

    Merge pull request #20 from DHBW-SE-2023/test_template
    
    added Test template

[33mcommit 5ff06773e18a89a9c25552085bbd98f645277407[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Thu Nov 23 14:22:32 2023 +0100

    added fyne test example

[33mcommit 6bc17b54931b798dc9a7857d2b70cb66e1afa3c4[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Thu Nov 23 14:20:22 2023 +0100

    update README

[33mcommit d6f716e1aecdce551e5cff8b066098cebc23a1c1[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Thu Nov 23 12:55:34 2023 +0100

    WIP Progress

[33mcommit 063998e2540192e998a49e0cc82e118582222da0[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Thu Nov 23 11:40:40 2023 +0100

    WIP

[33mcommit 08492c3d557e97e69c64e84bd2ec44e651f7de0f[m
Merge: e8de279 595025a
Author: Jonas Zagst <87520165+JonasZagst@users.noreply.github.com>
Date:   Thu Nov 23 10:06:18 2023 +0100

    Merge pull request #15 from DHBW-SE-2023/project-structure
    
    Setup project structure

[33mcommit de415443d37dda1ff0c88fde95a15545735a0537[m
Author: Eva <114435451+EvaMatz@users.noreply.github.com>
Date:   Thu Nov 23 09:31:54 2023 +0100

    Create LICENCE

[33mcommit d3db7e8ca2ee57c355b16b4767eae6b49eae9803[m
Author: danielSiegt08 <124031951+danielSiegt08@users.noreply.github.com>
Date:   Wed Nov 22 19:05:15 2023 +0100

    Added Simple Build/Test Pipeline

[33mcommit 595025a6b65491600822761db2eadde05ffc534b[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Wed Nov 22 15:25:56 2023 +0100

    added test example template

[33mcommit eb874080f59c9eab8dbbb51986cf68059cee3e99[m
Author: Vinzent <vinzent.engel@gmx.de>
Date:   Wed Nov 22 15:25:17 2023 +0100

    update readme

[33mcommit 6c54176b75a91e95d33398a490ddd09515e8fb11[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Wed Nov 22 10:40:43 2023 +0100

    Added templates

[33mcommit 9bbb291ea94f6b7916896c939c74a5112e982270[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 21 16:55:01 2023 +0100

    Moved global package variables to structs

[33mcommit 41a4187c45477fd6876598122b0381e359729e73[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 21 16:02:52 2023 +0100

    Updated backend layout

[33mcommit 9508e47d9984d88e723df51b80318843222c79e3[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Tue Nov 21 15:40:26 2023 +0100

    Refactored frontend layout

[33mcommit b7e49f451edb5e97935a76da032a8a1c8148484e[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Mon Nov 20 13:53:21 2023 +0100

    Updated go pkg link

[33mcommit e8de27992a06684bed83ba7d3fefd678d72df866[m
Merge: 7d8cd8d 2bc9fc2
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Fri Nov 17 16:28:27 2023 +0100

    Merge pull request #14 from DHBW-SE-2023/master
    
    Update dev to latest

[33mcommit 2bc9fc2d12febca9fef13bf068d81a8632e600c1[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Fri Nov 17 12:33:19 2023 +0100

    Create SECURITY.md

[33mcommit b585a51d6f0a50022e8b32baf2acb63dd3e53690[m
Merge: a607436 d2f0962
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Fri Nov 17 11:47:58 2023 +0100

    Merge pull request #9 from DHBW-SE-2023/dependabot/go_modules/golang.org/x/net-0.17.0
    
    Bump golang.org/x/net from 0.14.0 to 0.17.0

[33mcommit d2f0962593f40c5b2c0e3dde5a7016eff554138c[m
Author: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>
Date:   Fri Nov 17 10:44:38 2023 +0000

    Bump golang.org/x/net from 0.14.0 to 0.17.0
    
    Bumps [golang.org/x/net](https://github.com/golang/net) from 0.14.0 to 0.17.0.
    - [Commits](https://github.com/golang/net/compare/v0.14.0...v0.17.0)
    
    ---
    updated-dependencies:
    - dependency-name: golang.org/x/net
      dependency-type: indirect
    ...
    
    Signed-off-by: dependabot[bot] <support@github.com>

[33mcommit a607436562ebf30189a3a6ad8ec85c2c44a305a5[m
Author: MaxAlberti <maxalberti@gmx.de>
Date:   Fri Nov 17 11:37:08 2023 +0100

    Merged yaac-go-prototype into YAAC

[33mcommit 7d8cd8dec4de8417e002dfcce8a11a8d9a450d98[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Mon Oct 30 09:21:36 2023 +0100

    Update README.md

[33mcommit f5519a2fe5e67b6bc2f0944a6a4fdc6b2247470b[m
Author: Max <73841659+MaxAlberti@users.noreply.github.com>
Date:   Fri Oct 20 10:58:51 2023 +0200

    Update README.md

[33mcommit e49ea06bbb4d80c403fa234c0055a096b79a2e12[m
Author: Leander Gmeiner <leander.gmeiner@gmail.com>
Date:   Tue Oct 17 11:18:19 2023 +0200

    Initial commit
