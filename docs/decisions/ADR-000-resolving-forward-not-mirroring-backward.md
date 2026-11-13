# ADR-000: Resolving forward from ci-testing, not mirroring cloud-deploy backward

**Date:** 2026-07-04
**Status:** Accepted

## Context
Early in the reconstruction of this repo, inconsistencies found in ci-testing were
resolved by checking how forum-app-cloud-deploy — the later, integrator repo of the
same series — had solved the analogous problem, and porting that resolution back.
This was efficient, but it imported the reasoning of a later stage of the series into
an earlier one, which undermines the portfolio's central claim: that each repo
demonstrates its own reasoning, and that the series shows deliberate progression, not
retrofitted consistency.

## Decision
Every design or documentation decision in ci-testing is resolved independently,
grounded in TP6's own scope (unit testing automation, Services layer) and in
established software engineering concepts about unit testing — not in how
cloud-deploy resolved the same case. Once a decision here is fully fundamented, it
propagates forward: first to forum-app-qa-pipeline, then to forum-app-cloud-deploy —
never the reverse. Corrections to cloud-deploy are deferred until both earlier repos
in the series are finalized, so cloud-deploy receives one coherent pass reflecting
mature criteria, rather than multiple partial passes reflecting criteria still in
development.

## Alternatives considered
- Continue mirroring cloud-deploy's resolutions: faster, but conflates the
  reasoning of a later stage with an earlier one and weakens the series'
  progression narrative.
- Fix inconsistencies ad hoc without documenting the methodology itself: loses the
  reasoning trace that is the actual point of the portfolio — the "how do you think"
  demonstration this series exists to provide.

## Consequences
- Some fixes already identified as applicable to cloud-deploy (e.g., translating
  Spanish test names, replacing ADR-008's now-outdated risk rationale) are recorded
  as pending, not executed, until ci-testing and qa-pipeline are complete.
- Each repo's ADRs must stand on their own justification; cross-repo references are
  used only to explain divergence, never as the justification itself.

## Amendment (2026-11-12): which repository "forum-app-cloud-deploy" means here

When this ADR was written (2026-07-04), `forum-app-cloud-deploy` was the
repository now archived as `forum-app-cloud-deploy-legacy`; the resolutions
and the `ADR-008` this ADR mentions are that repository's. The current
`forum-app-cloud-deploy` was created on 2026-08-03; its `ADR-008` is about
explicit port binding, not test names. The decision is unchanged.
