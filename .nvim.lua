local opts = { noremap = true, silent = true }

vim.keymap.set("n", "<leader>tb", function()
	vim.cmd("!go build -o ./dist/neo ./cmd/neo && ./dist/neo")
end, opts)
