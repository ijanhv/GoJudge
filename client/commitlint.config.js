module.exports = {
  extends: ["@commitlint/config-conventional"], // Add a comma here
  rules: {
    "body-case": [2, "always", "sentence-case"],
    "body-max-line-length": [1, "always", 72],
    "header-max-length": [1, "always", 52],
    "type-enum": [
      2,
      "always",
      [
        "build",
        "chore",
        "ci",
        "deprecate",
        "docs",
        "feat",
        "fix",
        "perf",
        "change",
        "remove",
        "revert",
        "security",
        "style",
        "test",
        "refactor",
      ],
    ],
  },
};
