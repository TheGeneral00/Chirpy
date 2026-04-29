
local cjson = require('cjson')
local http = require('socket.http')

local home = os.getenv("HOME")
local testuserPath = home .. '/workspace/github.com/TheGeneral00/PostalWarden/json/testuser.json'

local testuserFile = io.open(testuserPath, 'r')
if not testuserFile then
    error("File not found: " .. testuserPath)
end

-- read entire file
local contentString = testuserFile:read("*a")
testuserFile:close()

local testuser = cjson.decode(contentString)

print("Testusers contain the following fields for this test:\n")

local keys = {}
for key, _ in pairs(testuser[1]) do
    table.insert(keys, key)
end

print(table.concat(keys, ", "))












