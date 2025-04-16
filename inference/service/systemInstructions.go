package inferenceService

const (
	//Default Return format
	defaultReturnFormat = "Return Format (Always Adhere to These Rules):" +
		" Responses must be strictly plain text, suitable for direct display in a textbox." +
		" Never use markdown formatting (no \"*\", \"-\", \"#\", or \"---\")." +
		" Do not return a response as if you were responding based off of a question." +
		" All your queries will be aggregated and decorated from the client. Do not give any indication you were prompted iteratively." +
		" If this task description cannot be answered or there is a gross amount of information missing to answer the task description, just remove it."
	//Default Warnings
	defaultWarnings = "Warnings:" +
		" output needs to be plain text absolutely no markdown" +
		" Always rely exclusively on the provided transcript without assumptions or inference beyond clearly available context."
)
