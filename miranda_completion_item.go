package main

// completionItem remains as a UI payload type for the completion popup.
// Legacy LSP integration has been removed.
type completionItem struct {
	Label  string
	Insert string
	Detail string
	Doc    string
}
