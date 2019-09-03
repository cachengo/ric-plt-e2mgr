##############################################################################
#
#   Copyright (c) 2019 AT&T Intellectual Property.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
##############################################################################

*** Settings ***
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     Collections
Library     REST      ${url}


*** Test Cases ***

Prepare Ran in Connected connectionStatus
    Post Request setup node b x-2
    Integer     response status       200
    Sleep  1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    CONNECTED


Send Reset reqeust with no cause
    Set Headers     ${header}
    PUT    /v1/nodeb/test1/reset    ${resetcausejson}
    Integer  response status  204