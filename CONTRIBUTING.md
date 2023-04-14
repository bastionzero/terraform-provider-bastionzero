# Contributing

This document details how to make contributions to this repo.

## Changelog

We follow Terraform provider plugin [changelog specifications](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#changelog-specification).

### Changie Automation Tool

The BastionZero Terraform provider uses the [Changie](https://changie.dev/)
automation tool for changelog automation.

To add a new entry to the `CHANGELOG`, install Changie using the following [instructions](https://changie.dev/guide/installation/)
and run:

```bash
changie new
```

then choose a `kind` of change corresponding to the categories specified [here](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#categorization).
Make sure to fill out the `body` following the entry format described below
(note: you can always update the auto-generated file later if you prefer to edit
the body in your text editor of choice).

Changie will then prompt for a Github issue or pull request number (note: If
your change spans across multiple issues or PRs, you can include all of them as
a comma separated list of numbers). _Repeat_ this process for any additional
changes. The `.yaml` files created in the `.changes/unreleased` folder should be
pushed to the repository along with any code changes.

#### Entry (`body`) format

Entries that are specific to _resources_ or _data sources_ should look like:

```markdown
resource/RESOURCE_NAME: ENTRY DESCRIPTION 

_or_

data-source/DATA-SOURCE_NAME: ENTRY DESCRIPTION
```

Do not include a trailing period as the generated file includes one for you.

#### Pull Request Types to `CHANGELOG`

The `CHANGELOG` is intended to show consumer-impacting changes to the codebase
for a particular version. If every change or commit to the code resulted in an
entry, the `CHANGELOG` would become less useful for consumers. The lists below
are general _guidelines_ to decide whether a change should have an entry.

##### Changes that should not have a `CHANGELOG` entry

* New tests or changes to existing tests
* Code refactoring

##### Changes that should have a `CHANGELOG` entry

* Major features
* Bug fixes
* Enhancements
* Deprecations
* Breaking changes and removals
* Dependency updates

##### **Special case**: Provider documentation updates

Typically, *do not* include a changelog entry if you are updating the provider
documentation in a PR that has other code changes as well. Only include an
entry, namely a `NOTES` entry that says something along the lines of "Update
docs", if your PR solely updates documentation.

This is a special case because noting a doc update can be noisy in the
CHANGELOG. It really only needs to be included if you intend to push the
documentation update immediately to the registry and there are no other
unreleased changes (because you can't make a release PR without having pending
unreleased changes).

## Releasing

Releasing a new version of the BastionZero Terraform provider is a
semi-automated process.

Use the ["Generate release pull
request"](https://github.com/bastionzero/terraform-provider-bastionzero/actions/workflows/gen-release-pr.yml)
workflow to auto-generate a release PR that collates all unreleased changes in
`.changes/unreleased` and updates the `CHANGELOG.md` accordingly. Specify `auto`
to let `changie` figure out what version to bump the provider to based on the
pending changes to be released. Otherwise, if you need finer control, specify an
explicit version number as follows: `v#.#.#`.

Please double check the auto-generated PR and squash merge the PR (do not
rebase).

When the release PR is merged to `master`, a tag pointed to the PR/commit that
updates the `CHANGELOG.md` is created and pushed to the repo. Furthermore, a
_draft_ release (pointing to the new tag) is created by
[`GoReleaser`](https://goreleaser.com/) containing build assets and release
notes parsed from the `CHANGELOG.md`. Please double check that the release looks
good and publish it when ready. When a release is published, an event is pushed
to a webhook which notifies the Terraform registry to ingest the new version and
make it available to the public.