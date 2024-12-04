---
title: "Help your code reviewer"
date: "2024-11-15"
description: "I've been reviewing PRs for years. Let me tell you what I want to see from a contributor."
categories:
  - "programming"
published: true
---

I've been reviewing pull requests (PRs) for years. Let me share what I want to
see from contributors, why some people consider code reviewing a waste of time,
and why I disagree.

In short:

**Don't waste your reviewer's time.**

Be diligent. Don't be lazy. Put yourself in the reviewer's shoes.

Open pull requests that you would like to review.


## Code reviews: a waste of time?

At different companies, I've noticed a code review flow that often looks like this:

1. The repository enforces a minimum number of reviewers (usually one or two).
1. The contributor submits a PR.
1. The contributor pings the reviewers in Slack.
1. The reviewers skim through the changes in the “Files changed” tab.
1. The reviewers approve the PR.

This process takes five seconds and, honestly, it's pointless. It's usually the
same one or two people taking the time to approve PRs instead of the
responsibility being shared across the team.


## PRs and code reviews in open source

If you look at popular open source projects, you'll notice that opening a PR is
never straightforward. Sure, you can open one, but you need to follow the
rules.

These rules [enforce consistency](./consistency.md), leading to higher-quality
code that stands the test of time.

I challenge you to try opening a PR to a big project...

- without having discussed the change with the team first, usually by opening
an issue
- providing no description ("No description provided", ugh...)
- with 38 commits, most of them being "fix" and "wip"
- changing 74 files
- fixing typos, refactoring code, and making unrelated changes that don't
depend on the feature you're working on.

_(don't actually do this 😱, it's a joke)_

**Now ask yourself this**: How can an open source project, managed
asynchronously over email or GitHub, coordinating the work of people without
weeks of Agile standups and retrospectives, be more organized than the average
company?

I see two reasons for this:

1. Ignorance.
2. Laziness.

And I don't mean that in a bad way. I'm here to help you improve your
PR-opening skills.


## How to help your code reviewer

First of all, read this: https://github.blog/developer-skills/github/how-to-write-the-perfect-pull-request/

Then, read this: https://codeinthehole.com/tips/advanced-pull-request-crafting/

Seriously, come back to my blog after you've read those two links. They cover
everything you need to write a PR that reviewers will love.

Still here? Alright, here's the gist:

1. Clearly describe the purpose of the PR and why the project needs this change.
Convince me it's a good idea.
1. Keep the PR small: a few commits are enough.
1. Make each commit atomic—each one should address one thing and pass your CI tests.
1. Tell a story with your commits. Show your thought process.
1. Review your own PR first. Double-check for mistakes, like accidentally
   committing your `node_modules` folder.

I'll repeat the beginning of this article:

> **Don't waste your reviewer's time.**


## A pragmatic tutorial

I spend a lot of time writing PRs, so here's an overview of my personal
workflow for rewriting history without overcomplicating things.

Let's say I want to add a new feature to a repo.


### Writing the first version

1. I create a branch from the main branch, I'll be the only one pushing to it:
    ```sh
    git switch -c my-feature
    ```
1. I write the code, frequently committing changes into "wip" commits:
    ```sh
    git add .
    git commit -m "wip"
    ```
1. Push it often. You don't want to lose your work if your laptop breaks.

### Cleaning up history

When it's time to open a PR, my messages will tell a story about what I did.
Let's clean up all the "wip" commits:

```
commit a146ef... (HEAD -> feature)
Author: Antonio
Date:   Today

    wip

commit 0410f...
Author: Antonio
Date:   Today

    wip

commit 0e52c...
Author: Antonio
Date:   Today

    wip

commit c1595... (main)
Author: Antonio
Date:   Yesterday

    first commit
```

Given that `c1595` is the latest commit I based my work on, I can reset my
branch to it:

```
git reset c1595
```

This will not change any of your local files, it will only delete the wip
commits, leaving all your work unstaged and ready to be committed again.

Now:

- you'll see all your changes in a single view (e.g. with `git diff`)
- you'll only see the final version of the files, their wip history is gone


### Crafting a story

Use tools like lazygit or your editor to stage specific changes. Aim for:

- Atomic commits that address one thing.
- "Buildable" commits that work on their own.

If your PR is small enough, a single commit is fine. But if you end up with
more than 7–8 commits, consider splitting the PR or explaining why it's not
possible in the description.

Let's spend some time crafting our PR history, writing a good message is out of
the scope of this article but a good read is
https://dhwthompson.com/2019/my-favourite-git-commit.

Having atomic commits helps when investigating a bug, `git bisect` is the best
tool to find the commit that introduced it. If we allow big, unbuildable
commits in our repo, pinpointing the exact problem will be harder.


### Opening the PR

Open the PR, optionally as a draft if it needs more work. Share it with your
team for feedback.

If reviewers suggest changes, avoid appending a "fix comments" commit. Instead,
rewrite the history to make it look like it was perfect the first time, the
following sections explain how to achieve this.


### Addressing reviewer feedback

If a reviewer asks to change something, refrain from appending a "fix comments"
commit to your PR.

> Disclaimer: this is an art that takes some experimentation to get right. I
> still sometimes end up frustrated and wished I could just append one commit
> to the end, still it gets easier with time.

As an example, let's say you have two commits, one introducing a new package
"foo" and the second changing some existing code to use the new package:

```
commit b45d6... (HEAD -> feature)
Author: Antonio
Date:   Today

    feat: use foo instead of bar in client contexts

commit 27a46...
Author: Antonio
Date:   Today

    feat: introduced foo

commit c1595... (main)
Author: Antonio
Date:   Yesterday

    first commit
```

If the reviewer asked to change something inside the "foo" package, you want to
change the `27a4` commit.

I achieve this using fixup commits. First of all, I'll make the changes I want,
then I'll create a new commit with the changes I made:

```
git commit --fixup 27a46
```

This will create a new commit, with a special message:

```
commit 1ca93... (HEAD -> feature)
Author: Antonio
Date:   Today

    fixup! feat: introduced foo

commit b45d6...
Author: Antonio
Date:   Today

    feat: use foo instead of bar in client contexts

commit 27a46...
Author: Antonio
Date:   Today

    feat: introduced foo

commit c1595... (main)
Author: Antonio
Date:   Yesterday

    first commit
```

You can squash the commits together with:

```
git rebase --interactive --autosquash main
```

Sometimes, you'll end up with conflicts. This is nothing to be scared of, it's
a normal part of the process.

If you look at the conflicts and see that something is wrong, you can always go
back by aborting the rebase:

```
git rebase --abort
```

Otherwise, fix the conflicts as usual and continue the rebase:

```
git rebase --continue
```

The advantage of fixup commits are:

- I can push them to the PR before squashing them, if I want to show them to
  the reviewer to confirm that I addressed their feedback
- I can batch multiple fixups, targeting different commits, and squash all of
  them in a single rebase


It's finally time to push your updates:

```
git push --force-with-lease
```

Don't forget to reply to the reviewer's comments, to make sure they're aware of
the changes.


### Rinse and repeat

Keep iterating until the reviewers are happy.

Remember: **your feature isn't done until the PR is merged**. It's your
responsibility to see it through.


## Conclusion

Opening good PRs and code reviewing are two sides of the same coin. I am
convinced that they are a crucial part of being a good software engineer.

Here are a few aliases I use multiple times a day:

```bash
# make a fixup commit, usage: `fixup 27a46`
alias fixup=git commit --fixup

# rebase interactively, usage: `rb origin/master`
alias rb=git rebase --interactive --autosquash

# rebase interactively on main, usage: `rbo`
alias rbo=rb origin/master

# make a fixup commit, then squash it immediately, usage: `fxrb 27a46`
# arguably the one I use the most
fxrb() {
    git commit --fixup "$1" && \
    rb "$1"^
}

# add staged changes to the last commit, essentially like a quicker fixup if
# the target commit is the last one, usage: `amndn`
alias amndn=git commit --amend --no-edit

# delete last commit, keep the changes, usage: `grh`
alias grh=git reset HEAD^

# git remote update, usage: `gru`
alias gru=git remote update

# git status
alias gs=git status
```

---

If you want to get in touch, find me on X:
[@zaphodias](https://x.com/zaphodias) or on Bluesky:
[@anto.pt](https://bsky.app/profile/anto.pt).
