## Server Architecture Overview

Kumo is built on a **Server-Driven UI (SDUI)** paradigm, where the server dictates the structure, layout, and visual presentation of the client applications. Rather than sending raw domain data for the client to parse and compute, the API acts as a product layout engine—delivering rich presentation data, styling variables, and explicit action directives.

### System Architecture Diagram

```text
                        +-----------------------------------+
                        |       Web Client (Frontend)       |
                        | - Renders Generic Widgets         |
                        | - Injects CSS Variables           |
                        | - Paper Doll / Avatar Compositor  |
                        +-----------------------------------+
                                         |  ^
                (1) REST/GraphQL Request |  | (2) JSON UI Payload & Data
                                         v  |
                        +-----------------------------------+
                        |         API Gateway / Auth        |
                        +-----------------------------------+
                                         |
                                         v
+-----------------------------------------------------------------------------------+
|                            Core Backend Service (Go)                              |
|                                                                                   |
|  +--------------------+     +--------------------+     +-----------------------+  |
|  |    User Engine     |     |   Goal/Plan Gen    |     |  SDUI Payload Builder |  |
|  | (Profiles, Auth)   |     | (Progress Tracker) |     | (Maps Goal -> Widgets)|  |
|  +--------------------+     +--------------------+     +-----------------------+  |
+-----------------------------------------------------------------------------------+
        |                                |                                 |
        | (Read/Write)                   | (Query)                         | (Trigger AI Intent)
        v                                v                                 v
+----------------+              +-----------------+             =====================
|  PostgreSQL DB |              | Redis / Cache   |             |    PROMPT HUB     |
| (User State,   |              | (Session State, |             | (Internal Pkg or  |
|  JSON Themes,  |<------------>|  Semantic Cache)|             |  Docker Gateway)  |
|  Schedules)    |              +-----------------+             =====================
+----------------+                                                         |
                                                                           | 1. Fetches Versioned Prompt
                                                                           | 2. Merges Variables
                                                                           | 3. Injects Vault Secrets
                                                                           | 4. Evaluates Cost Router
                                                                           v
                                                      +-----------------------------------------+
                                                      |  Third-Party AI Infrastructure Providers|
                                                      |                                         |
                                                      |  +-----------------+  +---------------+  |
                                                      |  | Tier 1: Fast    |  | Tier 2: Smart |  |
                                                      |  | (Gemini Flash / |  | (Reasoning /  |  |
                                                      |  |  DeepSeek)      |  |  Heavy Logic) |  |
                                                      |  +-----------------+  +---------------+  |
                                                      |          |                    |          |
                                                      |          +----------+---------+          |
                                                      |                     |                    |
                                                      |                     v                    |
                                                      |         +-----------------------+        |
                                                      |         |   Image Gen API       |        |
                                                      |         |   (Pixel Art Avatars) |        |
                                                      |         +-----------------------+        |
                                                      +-----------------------------------------+

                                                      Core Infrastructure Components
Core Backend Service (Go): The central business logic engine. It manages traditional application domains (user sessions, profile tracking, and timeline calculation) and transforms domain logic into structural UI components via the SDUI payload builder.

Prompt Hub (AI Gateway): A decoupled AI orchestration layer (implemented via a local package or a standalone gateway like Bifrost/BricksLLM). It abstracts third-party AI dependencies by handling prompt versioning, templating, and secure runtime injection of provider secrets.

Data Layer (PostgreSQL & Redis):

PostgreSQL: Acts as the persistent source of truth, storing user states, generated schedules, and historical layout configurations (eliminating the need to re-generate structural UI blocks on subsequent visits).

Redis: Powers the application's high-speed caching mechanics, managing active sessions and providing a semantic caching layer for standard conversational queries.

Key Architectural Patterns
Generic Descriptor Model: Visual elements are declared in the backend via generic design tokens (e.g., chaos.stat_card.v1). The server maps raw metrics into presentation-ready text fields before shipping, isolating the frontend client from schema modifications when layout content changes.

Action Modeling: Interactivity is server-driven. Component interactions include explicit action payloads (e.g., chaos.submit_log.v1). When a client event triggers, the payload is sent back to a universal endpoint, enabling the server to alter application state and return the modified layout dynamically.

Tiered Model Routing: To optimize API token spend, the Prompt Hub splits execution paths based on complexity. Routine operations (daily greetings, logging feedback) route to high-speed, cost-effective Tier 1 models. Heavy curriculum construction or critical analytical thinking tasks route to high-capability Tier 2 models.

Asynchronous Lifecycles: Resource-heavy workloads—such as generating custom 16-bit mentor avatar assets or compiling long-term text summaries—run outside the main execution block asynchronously, ensuring UI responsiveness remains isolated from AI processing latencies.
