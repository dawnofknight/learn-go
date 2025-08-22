-- Quick test to verify API key configuration
local handle = io.popen("printenv OPENAI_API_KEY")
local api_key = handle:read("*a"):gsub("\n", "")
handle:close()

if api_key and #api_key > 0 then
    print("✅ API Key found: " .. api_key:sub(1, 20) .. "...")
    print("✅ Length: " .. #api_key .. " characters")
    if api_key:match("^sk%-") then
        print("✅ Format looks correct (starts with 'sk-')")
    else
        print("❌ Format might be incorrect (should start with 'sk-')")
    end
else
    print("❌ API Key not found")
end
