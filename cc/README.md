# CC

Conventional commit cli

Enable consistent commit formatting in the style of **conventional commits** 

1. https://github.com/angular/angular/blob/22b96b9/CONTRIBUTING.md#-commit-message-guidelin
2. https://www.conventionalcommits.org/

### Parameters

- signed: the commit will be signed <https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits>
- breaking: the commit message will start with: type! or type(scope)!

### Commit Message Format

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

- *body and footer are optional*

### Type
0.  build:      Changes that affect the build system or external dependencies
1.  ci:         Changes to our CI configuration files and scripts
2.  docs:       Documentation only changes
3.  feat:       A new feature
4.  fix:        A bug fix
5.  refactor:   A code change that neither fixes a bug nor adds a feature
6.  test:       Adding missing tests or correcting existing tests
7.  chore:      Ad-hoc task that doesn't match other types

### Scope

The scope should be the name of the smallest unit of code that has been affected.

### Subject
The subject contains a succinct description of the change:

### Body
The body should include the motivation for the change and contrast this with previous behavior.

### Footer
The footer should contain any information about Breaking Changes and is also the place to reference GitHub issues that this commit Closes.

Breaking Changes should start with the word BREAKING CHANGE: with a space or two newlines. The rest of the commit message is then used for this.

### Writing Guidlines
- use the imperative, present tense: "change" not "changed" nor "changes"
- don't capitalize the first letter
- no dot (.) at the end
