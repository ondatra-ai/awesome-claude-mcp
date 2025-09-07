# Q&A on Requirements for the MCP Server for Editing Google Docs

**Question 1:** In which interface do you expect to use this tool most often — in the cloud assistant’s web interface or in the desktop app?

**Answer:** Yes, I want to use both, but mostly the Claude web interface.

**Question 2:** When you use this tool through the web interface, do you expect edits in Google Docs to be applied automatically, or should the tool first show you a draft/confirmation and only then apply them?

**Answer:** This tool simply edits the document, that’s it. Claude decides everything else.

**Question 3:** Do you want only you personally to use this tool, or do you allow other people (e.g., colleagues) to edit their documents through it as well?

**Answer:** I’d like everyone in my organization to use the tool with their own credentials. They can edit only their own documents.

**Question 4:** How do you envision selecting a document to edit — will users paste the link/ID manually, or should the tool allow searching and choosing documents?

**Answer:** The document ID is provided to this tool by another tool that searches Google Drive.

**Question 5:** Should each user sign into their Google account once and then the tool works automatically, or do you prefer confirming the login every time?

**Answer:** Sign in once.

**Question 6:** Should we limit it to text Google Docs only, or should the tool handle other formats too?

**Answer:** Text documents only.

**Question 7:** Should the tool be able to leave comments and suggestions in documents, or only make direct edits?

**Answer:** Direct edits only.

**Question 8:** When the tool edits text, what’s more important to you: speed or reliability?

**Answer:** I don’t care about speed.

**Question 9:** Should the tool keep an activity log?

**Answer:** No, nothing additional is needed.

**Question 10:** Should the tool be triggered only manually via commands in Claude, or also automatically by events/triggers?

**Answer:** Only via Claude when Claude calls it through the MCP protocol.

**Question 11:** Should the tool be available only inside the organization, or can it be a public service with OAuth protection?

**Answer:** It should be a shared service. Anyone can use it via a link with OAuth protection.

**Question 12:** Should changes affect formatting and structure, or only text?

**Answer:** The tool receives Markdown and formats the document according to the Markdown structure.

**Question 13:** Should Markdown completely replace the content, or can it also insert/edit?

**Answer:** It must allow inserting, editing, and fully overwriting the document.

**Question 14:** How do we specify the edit location in the document?

**Answer:** Via an anchor or a regular expression. You can specify to replace the match, insert before it, or after it.

**Question 15:** If there are multiple matches, what is the default behavior?

**Answer:** Replace all matches.

**Question 16:** What should happen if the anchor isn’t found?

**Answer:** Return an error with options for Claude on how to proceed.

**Question 17:** If the input Markdown contains images, what should we do?

**Answer:** Insert the images into the document.

**Question 18:** Should we insert hyperlinks from Markdown?

**Answer:** Yes, where possible.

**Question 19:** Should we preserve the heading structure from Markdown?

**Answer:** Yes, convert them to Google Docs headings.

**Question 20:** What should the tool return?

**Answer:** Structured JSON with the result so Claude can show a report or next‑step options.

**Question 21:** Case sensitivity when searching for anchors?

**Answer:** Ignore case (case‑insensitive).

**Question 22:** Should we provide protection against document version conflicts?

**Answer:** No, the tool simply edits documents without checks.

**Question 23:** Should the tool be able to create new documents or only update existing ones?

**Answer:** Only update existing documents by ID.

**Question 24:** If the Google API is unavailable, should we retry or return an error?

**Answer:** Just return an error to Claude.

**Question 25:** Which image sources should be supported?

**Answer:** Only external URLs accessible without authorization.

**Question 26:** If the Markdown contains tables, how should they be inserted?

**Answer:** Convert them into Google Docs tables.

**Question 27:** Do we need a limit on the size of the input Markdown?

**Answer:** No limits.

**Question 28:** If only `docId` and Markdown are passed without a mode, what should happen?

**Answer:** Overwrite the entire document by default.


# MCP Tool JSON

```json
{
  "name": "google_docs_editor",
  "description": "Edits existing Google Docs by ID using Markdown content.",
  "input_schema": {
    "type": "object",
    "properties": {
      "docId": {
        "type": "string",
        "description": "Google Docs document ID"
      },
      "markdown": {
        "type": "string",
        "description": "Markdown text to insert or replace"
      },
      "mode": {
        "type": "string",
        "enum": [
          "replace_all",
          "append",
          "prepend",
          "replace_match",
          "insert_before",
          "insert_after"
        ],
        "description": "Edit mode"
      },
      "anchor": {
        "type": "string",
        "description": "Text or regular expression used to locate the insertion point"
      },
      "case_sensitive": {
        "type": "boolean",
        "default": false,
        "description": "Whether to respect case when searching for the anchor"
      }
    },
    "required": ["docId", "markdown"]
  },
  "output_schema": {
    "type": "object",
    "properties": {
      "type": {
        "type": "string",
        "enum": ["ok", "error"]
      },
      "docId": {
        "type": "string"
      },
      "matches_found": {
        "type": "integer"
      },
      "matches_changed": {
        "type": "integer"
      },
      "preview_url": {
        "type": "string"
      },
      "warnings": {
        "type": "array",
        "items": { "type": "string" }
      },
      "message": {
        "type": "string"
      },
      "hints": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "action": { "type": "string" },
            "label": { "type": "string" }
          }
        }
      }
    },
    "required": ["type"]
  }
}
```


# Examples of Using the MCP Tool

## Example 1. Full document replacement

**Request:**
```json
{
  "name": "google_docs_editor",
  "arguments": {
    "docId": "1AbCdEfGhIjKlMnOpQrStUvWxYz123456",
    "markdown": "# New report\n\nThis is the updated content of the document.",
    "mode": "replace_all"
  }
}
```

**Response:**
```json
{
  "type": "ok",
  "docId": "1AbCdEfGhIjKlMnOpQrStUvWxYz123456",
  "matches_found": 1,
  "matches_changed": 1,
  "preview_url": "https://docs.google.com/document/d/1AbCdEfGhIjKlMnOpQrStUvWxYz123456",
  "warnings": []
}
```

---

## Example 2. Insert text before the anchor

**Request:**
```json
{
  "name": "google_docs_editor",
  "arguments": {
    "docId": "1AbCdEfGhIjKlMnOpQrStUvWxYz123456",
    "markdown": "## New section\n\nText of the new section.",
    "mode": "insert_before",
    "anchor": "replace me"
  }
}
```

**Response:**
```json
{
  "type": "ok",
  "docId": "1AbCdEfGhIjKlMnOpQrStUvWxYz123456",
  "matches_found": 2,
  "matches_changed": 2,
  "preview_url": "https://docs.google.com/document/d/1AbCdEfGhIjKlMnOpQrStUvWxYz123456",
  "warnings": []
}
```

---

## Example 3. Error: anchor not found

**Request:**
```json
{
  "name": "google_docs_editor",
  "arguments": {
    "docId": "1AbCdEfGhIjKlMnOpQrStUvWxYz123456",
    "markdown": "Text to insert.",
    "mode": "replace_match",
    "anchor": "PLACEHOLDER"
  }
}
```

**Response:**
```json
{
  "type": "error",
  "code": "ANCHOR_NOT_FOUND",
  "message": "Anchor 'PLACEHOLDER' not found in the document.",
  "hints": [
    { "action": "insert_at_end", "label": "Insert at the end" },
    { "action": "replace_all", "label": "Overwrite the entire document" },
    { "action": "ask_user", "label": "Ask the user" }
  ]
}
```
