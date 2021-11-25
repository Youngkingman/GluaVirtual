local function max(...)
    local args = {...}
    local val,idx
    for i = 1, #args do
        if val == nil or args[i] > val then
            val,idx = args[i],i
        end
    end
    return val, idx
end

local function assert(v)
    
end

local v1 = max(3,4,5,6,6,7,4,2,15,67)
assert(v1 == 67)
local v2, i2 = max(3,9,7,128,35)
assert(v2 == 128 and i2 == 4)
local v3,i3 = max(max(3,9,7,128,35))
assert(v3 == 128 and i3 == 1)
local t = {max(3,9,7,128,35)}