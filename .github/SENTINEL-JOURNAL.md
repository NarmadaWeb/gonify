# SENTINEL SECURITY JOURNAL

## 2025-05-22 – Unbounded Response Buffering in Middleware
**Pattern discovered:** Minification middleware buffering entire response bodies without size limits.
**Business/project impact:** Enables Denial of Service (DoS) via Out-Of-Memory (OOM) crashes by any user capable of triggering a large response.
**Constraint / future rule:** All middleware that buffers response bodies MUST enforce a maximum size limit (MaxBodySize).
**Recommended controls:** Default limits (e.g., 10MB) and early skipping for encoded or oversized content.
