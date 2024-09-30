"use client";

import React  from "react";
import * as monaco from "monaco-editor";
import { Editor, loader } from "@monaco-editor/react";

// Initialize the Monaco loader
loader.config({ monaco });

interface MonacoEditorProps {
  value: string;
  onChange: (value: string) => void;
  language: string;
  theme?: "vs-dark" | "light";
}

const MonacoEditor: React.FC<MonacoEditorProps> = ({
  value,
  onChange,
  language,
  theme = "vs-dark",
}) => {


  return (
    <Editor
    height="90vh"
    defaultLanguage={language}
    defaultValue={value}
  />
  );
};

export default MonacoEditor;
