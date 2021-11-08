Tab1 = { key1 = "val1", key2 = "val2", "val3" ,"yinbaoheihei"}
for k, v in pairs(Tab1) do
    print(k .. " - " .. v)
end

Tab1.key1 = nil
for k, v in pairs(Tab1) do
    print(k .. " - " .. v)
end