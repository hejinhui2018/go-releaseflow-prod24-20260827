# ReleaseFlow

ReleaseFlow is a Go command-line release-packet tracking service for internal software delivery teams. It records submit, review, approval, publish, and rollback actions for each release packet, then serves the recovered state and audit history to release coordinators after a restart.

The CLI is intended for release coordinators who need a compact local record of each packet, including reviewer notes and publication status. It also exports a machine-readable report for downstream handoff.
