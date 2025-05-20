# ✨ Contributing to Gonify ✨

Hello, Awesome Contributor! 👋

Thank you so much for your interest in contributing to **Gonify**! We're
thrilled to welcome any help, whether it's ideas, bug reports, code fixes,
documentation, or even just spreading the good word about Gonify. Every
contribution, no matter how small, is incredibly valuable to us and the
community.

This document is a guide to help you contribute. Let's make Gonify even more
awesome together! 🚀

## 📜 Code of Conduct

Before you start, please make sure you've read and agree to abide by our
[Code of Conduct](CODE_OF_CONDUCT.md). We are committed to maintaining a
friendly, inclusive, and respectful environment for all contributors.
(Ensure you create this `CODE_OF_CONDUCT.md` file, perhaps adapted from the
Contributor Covenant).

## 💡 How Can I Contribute?

There are many ways to contribute to Gonify:

* 🐞 **Reporting Bugs:** Found something not working as expected? Let us know!
* 🌟 **Suggesting New Features:** Have a brilliant idea to make Gonify better?
  We'd love to hear it!
* 🛠️ **Fixing Code:** Check out our
  [issue list](https://github.com/NarmadaWeb/gonify/issues) for bugs or
  features you can work on.
* 📚 **Improving Documentation:** Found a typo, unclear explanation, or a
  section that could be improved? Your help is greatly appreciated!
* 🗣️ **Spreading the Word:** Tell your friends or colleagues about Gonify!
* 🎨 **Improving UI/UX (if applicable):** If Gonify has visual aspects,
  suggestions for improvement are always welcome.

## 🚀 The Pull Request (PR) Process

We love Pull Requests! Here's a guide to help your PR get reviewed and merged
smoothly:

1. **Initial Discussion (Important!):**
    * For significant changes, new features, or major refactors,
      **please discuss it first** via our
      [Issue Tracker](https://github.com/NarmadaWeb/gonify/issues).
      Create a new issue if one doesn't already exist.
    * For minor bug fixes or typos, you can go ahead and create a PR, but be
      sure to reference any related issue if one exists.
    * We're also on `[Link to your Gonify Discord/Slack/Forum - if any]` for
      more casual discussions.

2. **Fork & Branch:**
    * **Fork** the `NarmadaWeb/gonify` repository to your GitHub account.
    * **Clone** your fork to your local machine:
      `git clone https://github.com/YOUR_USERNAME/gonify.git`
    * Create a **new branch** from `main` (or other primary development branch)
      for your work: `git checkout -b feature/your-cool-feature-name` or
      `fix/brief-bug-description`.

3. **Coding & Testing:**
    * Make your code changes.
    * **IMPORTANT:** If you're adding a new feature or fixing a bug,
      **you MUST include relevant unit tests**. We believe in well-tested code!
    * Ensure all tests pass: `go test ./...` (or your project's specific test
      command).
    * Ensure your code follows Go standards (`gofmt`, `golint`/`golangci-lint`
      if used).

4. **Awesome Commit Messages:**
    We use emoji prefixes for commit messages to make them easier to read and
    understand. Choose the one that best fits:

    > 🔥 `Feature`: Adding a cool new functionality.
    > 🩹 `Fix`: Fixing a bug or issue.
    > ♻️ `Refactor`: Tidying up code without changing functionality.
    > 📚 `Doc`: Changes to documentation (README, code comments, etc.).
    > 🎨 `Style`: Fixes related to formatting, whitespace, without logic
    > changes.
    > 🚨 `Test`: Adding or improving tests.
    > ✨ `Chore`: Other tasks like dependency updates, build scripts, etc.
    > 🚀 `Release`: Used when preparing a new release (usually by maintainers).

    **Commit Message Examples:**
    * `🔥 Feature: Implement email notification module`
    * `🩹 Fix: Address panic on empty input in X function`
    * `♻️ Refactor: Simplify user validation logic`
    * `📚 Doc: Update README with new usage examples`
    * `🚨 Test: Add test case for edge scenario in Y module`

5. **Create a Pull Request:**
    * Push your branch to your fork on GitHub:
      `git push origin feature/your-cool-feature-name`
    * Open the `NarmadaWeb/gonify` repository on GitHub, and you'll see a
      prompt to create a Pull Request from your branch.
    * Provide a clear title and a detailed description for your PR. Explain
      **what** you changed and **why**.
    * If your PR addresses a specific issue, don't forget to link it in the
      description (e.g., `Closes #123`).
    * Our team will review your PR. There might be discussions or requests for
      changes. We'll try to respond as quickly as possible!

## 🔧 Setting Up Your Development Environment (Optional, but Helpful)

To start contributing code to Gonify:

1. Ensure you have **Go version X.Y.Z** or newer installed
   ([https://golang.org/dl/](https://golang.org/dl/)).
2. Clone the repository as described in the "Fork & Branch" section.
3. Navigate to the project directory: `cd gonify`
4. Download dependencies (if using Go Modules): `go mod tidy` or `go get ./...`
5. Run tests to ensure everything is working: `go test ./...`
6. (Add any other specific steps if necessary, e.g., database setup,
   environment variables, etc.)

## 👍 Other Ways to Support Gonify

If you love Gonify and want to show your support (besides code contributions):

1. ⭐ Give our [GitHub repository](https://github.com/NarmadaWeb/gonify)
   a **Star**!
2. ✍️ Write a **blog post, review, or tutorial** about Gonify on platforms
   like Medium, Dev.to, or your personal blog.
3. 💬 Join the discussion on our
   [Issue Tracker](https://github.com/NarmadaWeb/gonify/issues)

---

Thank you for taking the time to read this contributing guide. We can't wait
to see your contributions! If you have any questions, please don't hesitate
to ask.

Happy contributing! 🎉
**The Gonify Team (NarmadaWeb)**
