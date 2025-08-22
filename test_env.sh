#!/bin/bash

echo "🧪 Testing Environment Variable Access"
echo "======================================"

# Test from shell
echo "1. Shell environment:"
echo "   OPENAI_API_KEY length: ${#OPENAI_API_KEY} characters"
echo "   First 20 chars: ${OPENAI_API_KEY:0:20}..."

# Test from Neovim
echo ""
echo "2. Neovim environment test:"
nvim --headless -c "
lua << EOF
local api_key = vim.fn.getenv('OPENAI_API_KEY')
if api_key and api_key ~= vim.NIL and #api_key > 0 then
    print('✅ Neovim can access API key')
    print('   Length: ' .. #api_key .. ' characters')
    print('   First 20 chars: ' .. api_key:sub(1, 20) .. '...')
else
    print('❌ Neovim cannot access API key')
    print('   Value: ' .. tostring(api_key))
end
EOF
qa!
"

echo ""
echo "3. Testing plugin loading:"
nvim --headless -c "
lua << EOF
local ok, chatgpt = pcall(require, 'chatgpt')
if ok then
    print('✅ ChatGPT plugin loads successfully')
else
    print('❌ ChatGPT plugin failed to load')
    print('   Error: ' .. tostring(chatgpt))
end
EOF
qa!
"
