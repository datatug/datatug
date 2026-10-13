# Context Bag product and UX proposal

2026-10-13 · Approved M0–M3 direction · Synthetic examples throughout

## The experience to aim for

“I have Order 123. Show me what happened.” Start a session, paste `123` into the
suggested **Order → ID** field, and Apply Context. Orders shows the order; opening
OrderDetails shows its lines. The context banner says why. Attach the entry to chat
and ask a question. Save a finding to a Detective Case. Stop for lunch. Tomorrow,
resume exactly this session, or start another without carrying yesterday's filters.

Context is global within the selected session across queries, Explorer, dashboards and
Checks. It is not restricted to a query editor, and is never ambient across users or sessions.

The product promise is **carry the question across tools without losing the evidence**.
The bag is a small workbench, not a prescribed workflow: an order, a date window, two
cohorts and an invoice number can coexist. Tags describe or address knowledge; bag
entries supply potential runtime inputs. Neither replaces the other.

## Object vocabulary and information architecture

| Object | Job | What it does not imply |
|---|---|---|
| Session | Resumable workspace with owner, source scope, bag, tabs and conversation | Confirmed problem or active process |
| Context Bag | One of a session’s environment-scoped collections of typed entries | Automatically attached to every chat |
| Entry | Stable identity, semantic subject, editable label, revision history | One string value forever |
| Applied context | Immutable manifest of enabled entry revisions and combination choices | Current head of every entry |
| Case (Detective) | Durable question, findings, evidence and closure across sessions | A second mutable copy of the session bag |
| Incident (Incidentius) | Confirmed problem, severity and operational lifecycle | Every investigation is an incident |
| Check | Reusable validation using explicit inputs and recorded execution | A session or investigation recipe |
| Scenario | Optional guide for investigating | Mandatory sequence or Check evaluator |

One session may contribute evidence to several Cases; one Case may collect several
sessions. Begin with an optional primary Case link and explicit additional evidence
links. “Start investigation from chat” can retain the session ID and change its
workspace mode; offer **Branch into new session** when independent work is intended.
An incident opens an overview with **Start investigation** and selected context preview.

Navigation: project → Sessions → session workspace → optional Case/evidence details →
back to the same session. The owner selected **editor-style session tabs at the top of
the context workspace** as the initial switcher. Session tabs sit above the compact
hybrid context summary; query/dashboard view tabs remain inside the selected session,
with a separate visual hierarchy. Cold routes carry only opaque navigation IDs;
the authorized store restores context, never URL values.

Each session tab shows its title and a restrained state/pending indicator; its full
accessible name identifies the session and state. New session and All sessions remain
visible beside a horizontally scrolling tab strip. All sessions provides search and
running/stopped/archived filters, reopening a session in the strip without resuming it.
Keep this list reusable for a possible ChatGPT-style persistent sidebar later; that
alternative is not part of the initial layout. Session rows identify database location;
tabs need not repeat every bag's environment because one session can contain several.

Closing a tab changes only this device's open-session list. It does not Stop, Archive,
Delete, discard a draft or cancel a run. Close active selects an adjacent open session;
closing the last tab shows the session chooser, without silently creating a new session.
Persist recoverable drafts before switching/closing; a save failure offers Retry or Stay.
Reopening restores the same bags, applied state, composer and views without executing.
Use roving focus with Left/Right/Home/End; Enter/Space activates the focused tab.
Close actions have separate accessible buttons, never buttons nested inside a tab button.
On narrow screens retain scrolling tabs plus All sessions; do not wrap multiple tab rows
above the editor. Opening a deep link also opens/focuses its session tab after authorization.

## Session lifecycle: activity is separate from conclusion

Recommended states are `running`, `stopped`, `archived`, plus a deletion tombstone.
Running means available for user work; it does not create background AI activity.
Multiple sessions can be running. `focusedSessionId` is a device UI preference.

| Action | Observable result | Guard / recovery |
|---|---|---|
| New | New stable session with empty bag, explicit database, draft composer | Creation failure leaves entered title; no fake saved state |
| Stop | Drafts persisted and still editable; no new execution; best-effort cancel owned work | Show still-running/unknown jobs until actual receipts settle |
| Start / Resume | Same context, tabs, drafts and history; execution becomes available | Never rerun prior queries just by resuming |
| Archive | Stop first, read-only workspace hidden from default list | Pending jobs must settle or remain visibly tracked; Restore → stopped |
| Delete | Impact preview, explicit confirmation, hidden tombstone | Block unresolved jobs; linked evidence retained or deletion blocked by policy |
| Switch | Focus changes; other session state stays intact | Dirty drafts save locally; failed save offers retry or stay |
| Close tab | Remove from this device's open tabs; session retains its lifecycle | Reopen through All sessions; preserve drafts and in-flight receipts |

Stopping a session does not close a Case. Closing a Case requires a reason and changes
the investigator's conclusion. Archiving a Case does not silently archive every linked
session. Deleting a session is not a promise to erase Git history, shared evidence or
backups; show exactly what is removed and retained. Permanent purge is a separate
capability, with honest store-specific limits.

Scope is `{database, principal/security context, project, environment, identityDomain}`.
A session has immutable project/database access scope. Each bag has an immutable
environment and identity domain. Use another bag for another environment within the
same project; changing project/database creates a reviewed session branch. No historical
manifest is reinterpreted and no production IDs silently become UAT IDs.
The saved session remains readable only under current authorization.

## Selection and revision state machines

Maintain **three independent states**: draft session selection, applied session
manifest, and composer attachments. This is worth the complexity internally because
it prevents two common errors: believing a checkbox has already changed a result,
and believing a historical answer is still about the current cohort.

```text
Entry edit → new immutable entry revision → “Update available”
Checkbox / master / combination edit → draft selection → “Changes not applied”
Apply → validate complete manifest → persist selection revision → refresh compatible Explorer views
Attach → composer reference chip → Send → resolve + pin revision → immutable turn input
```

1. Master checkbox and entry checkboxes only modify the draft. Master off retains the
   individual selections; Apply suspends automatic bag binding in all session environments across Explorer,
   query inputs, dashboards and Checks. Explicit manual inputs/constants and composer attachments remain independent. Master on then Apply restores
   them. “Ignore context / View all” stages master off and offers Apply, rather than
   introducing a second hidden switch. A per-view override, if added, must be explicit. Per-environment applied slots are
   independent: UAT entry edits cannot change PROD context. The explicitly labeled
   All session context master is the one global suspension gate.
2. Apply first validates the manifest and its chosen combination syntax. Consumer-specific
   ambiguity, missing input or unsupported shapes are reported per view, not a global
   veto. Manifest validation failure preserves the previous applied manifest. On success the manifest is committed once; each table
   independently shows loading/result/error tied to that manifest. No distributed
   all-or-nothing query execution is implied. Old rows are labeled stale while loading.
3. Entry edits, membership changes and parent refreshes never replace pinned revisions
   automatically. Show `Bad orders · 200 · applied r4 · r5 available` until Apply.
4. Unchecking everything means no active constraints after Apply. A checked empty
   collection means zero matches. Unknown/unavailable values mean an error, never all.
5. A selected entry irrelevant to a table is labeled “Not used here.” A required input
   that cannot bind blocks that view. A wholly unrelated table visibly says “No applicable
   context” and requires a clear View all choice before unconstrained initial loading.
6. While Apply runs, disable duplicate submit. A request captures session ID, applied
   revision, scope and UI generation. Responses append to their originating session;
   stale generations cannot overwrite the focused grid.

## One ID, cohorts and conditions

The first input uses current selected field/schema as a suggestion, explicitly labeled
**Order → ID**. With no trustworthy suggestion, choose an entity then field; do not
infer that every integer is an order. Semantic registry search narrows to this project's
mapped schema before broader suggestions. Missing semantics permits explicit physical
field binding, visibly local to that source, never a made-up canonical MeaningGraph ID.

Default address `contextBag.order.orderId` is a convenience alias, not persisted query
syntax. Stable opaque `entryId` survives rename and value replacement. A duplicate
semantic entry shows its field, count and creation time; the UI invites Good orders / Bad orders without
requiring an alias for the one-ID path. Aliases are unique within the bag; renames keep
reference identity. Typed equality distinguishes number 123 from text "00123".

Paste entry supports one-per-line and quoted CSV modes. Preview count, duplicates,
invalid values and type conversion before save. Pasting 5,300 IDs never renders 5,300
chips. Lists preserve order/duplicates only when declared `list`; identity cohorts use
`set` with exact typed deduplication. Null, missing and empty text are distinct states.

When two entries match one input, never choose the first:

| Intent | Recommended interaction | Result meaning |
|---|---|---|
| Inspect Good only | Check Good, uncheck Bad, Apply | One named input set |
| Inspect all suspects | Select **Union** in binding preview | Deduplicated union; overlaps counted once |
| Compare Good/Bad | Select **Compare** | Separate labeled result tabs by default |
| Intersect / subtract | Advanced explicit collection operation later | Named derived revision; no inference from labels |

Different compatible fields combine with explicit AND in an inspect view. Multiple
same-field scalar candidates also require a choice; scalar query inputs cannot accept
a union. “Good” and “Bad” are human labels, not proof of a disjoint partition. Show
overlap from the two resolved cohorts. Unclassified counts require an explicit parent universe.
Later classification can enforce mutual exclusion, with atomic moves and reasons;
free overlapping cohorts remain available.

Show both **100 input orders → 284 OrderDetails rows**. In Compare, results may overlap
and shared invoice lines may appear on both sides. A combined grid is an alternative,
with membership badges, unique row counts and a clear overlap explanation; do not
sum returned row counts to infer number of distinct entities.

Condition entries are typed predicates, not SQL fragments. Store `gt`, range or relative
time expressions as an AST. Only compatible declared input shapes accept them. The
first slice supports scalar/set; range next. “Last five hours” stores its expression,
timezone and the absolute half-open interval resolved at Apply/Send. A later live mode
must produce explicit new executions, not rewrite previous evidence.

## Explorer, query editor and chat have different promises

Explorer table opening is an intentional read action. It resolves applied context through
known MeaningGraph binding + ModelSpec keys/references + physical source mappings.
Direct Order → OrderDetails is the first supported path. Order → InvoiceLines may require
a bridge or multiple candidates; show relationship choices and estimated fanout. Similar
names alone are never evidence. Reject unsupported or policy-disallowed joins.

Saved query definitions, TugQL/TugQTree and Checks stay portable. Their input pane shows
`orderIds ← Bad orders r4` with override/removal and type validation. Explicit Run remains
required. Applying context refreshes compatible Explorer views; it does not launch all
open saved queries or Checks. Bound input provenance is shown before and after execution.

Every bag entry offers Attach to Composer. Header offers **Attach checked entries** and
**Attach all entries…** with a preview; neither is secretly affected by session master
off. Composer chips are removable and scoped to this session. Chips default to following the latest revision until Send, with a visible changed marker
and an optional Pin this revision action. Send resolves following references and pins all
inputs atomically; missing/incompatible references require correction. No routine confirmation
is added for a permitted latest revision. A retry uses the same send operation and snapshot.

Large attachment to the model: label, semantic identity, 5,300 members, revision,
allowed operations and opaque authorized handle. Values go to the executor through
runtime binding. A single ID may enter the prompt; samples require policy permission
and bounded disclosure. A model-supplied handle is untrusted and reauthorized at use.

Past chat results never change after Apply, cohort edits or parent refresh. “Filter this
result” / “make a chart” creates a derived view revision with parent result ID, recorded
transform and Undo. Rerun creates a new execution, since source data may have changed.
Pinned inputs make provenance reproducible; they do not promise the remote data is frozen.

## Extraction, evidence and knowledge

Flow: paste transient text → select local/external extraction → provider disclosure →
candidate review → correction → commit bag entries → separate evidence-retention choice.
The paste is not automatically a persisted ordinary chat turn. A normal already-sent
chat message cannot be retroactively described as unretained extraction input.

Example synthetic text: “Good orders: 101, 102. Suspect: 123, 124. Invoice INV-8.
Label ZX-9; maybe 125.” Review separates four semantic fields/cohorts, flags uncertain
125, and allows removal. Do not show invented numeric confidence; show supporting
snippet and uncertainty. Schema narrowing uses deterministic candidates first; Jev is
an integration candidate, not a prerequisite assumed deployed.

Initial source retention is Always ask: whole source / selected snippet / don't retain.
Explain that an external provider sees submitted text even when DataTug does not retain
it. Declining retention leaves only reviewed values, source-kind metadata and operation
receipt. Do not preserve raw text in audit, telemetry, crash reports or hidden prompt logs.

Annotations target an entity key, cohort revision, execution or explicit selected snapshot.
Origin separately records session, Case and incident references. Unified tag records hold
free labels or validated structured `entity:` addresses. Opening Order 123 shows **Direct
notes** separately from **Included in commented cohorts**. The latter needs actual indexed
membership in the referenced revision, not propagation of every cohort tag to every row.
Operational observations can propose semantic knowledge for review; they never mutate
MeaningGraph definitions automatically.

## Placement alternatives: use the same seven scenarios

The original prototype compared four layouts with shared fixtures. Owner review selected
hybrid and rejected top-only; left/right alternatives remain research references. Validate
the selected layout at laptop width, small screen and keyboard-only.

| Layout | Strength | Cost / failure risk | Best fit |
|---|---|---|---|
| Top context bar | One ID always visible; full grid width | Owner rejected this context layout; retain only as research history | Not a delivery option |
| Left sidebar | Bag adjacent to project/tree navigation | Competes with Explorer and session navigation; cramped grids | Investigation-first mode |
| Right panel | Full editor with provenance/revisions | Competes with row Inspect; long way from navigation | Deliberate cohort work |
| Hybrid: compact status strip + drawer | Always shows applied state; editor opens only when needed | Needs an obvious affordance; drawer can hide draft changes | Owner-selected direction |

Selected hybrid: compact summary **Context · 2 applied · 1 pending** opens a right drawer
that temporarily replaces Inspect, not a fourth permanent column. On narrow screens use
a full-height sheet with sticky Apply and return focus to the opener. In terminal use a
Context tab/pane and status line, with no printable-key commands while typing. The
cohort comparison display remains separately selectable: tabs first, combined view tested.

| Script in every layout | Observable success |
|---|---|
| One Order 123 | Useful Orders result in a few deliberate actions, no naming ceremony |
| 300 suspects: 100 good / 200 bad | Choose union or compare; distinguish input and output counts |
| 5,300 IDs → composer | Compact chips, immutable send receipt, no ID dump |
| Disable context | Master change remains pending; Apply suspends and retains selections |
| Attach selected | Uses checked entries, removable chips; no unexpected “all” |
| Open OrderDetails | Applied context follows verified relationship, without second Apply |
| Incident → investigations | Preview selected snapshot; children independently editable |

Also test switching with drafts and a running query; stopped session prevents execution; archive
restore; referenced deletion; empty set; ambiguous join; permission revoked; offline
source unavailable. Use visible text/icons as well as color, native checkbox labels,
logical tab order, announced pending/loading/error states and focus restoration.

Usability experiment (proposed, not performed): five investigators using the selected
hybrid layout; measure time to correct first result, wrong-scope executions, recognition
of pending versus applied, and ability to explain what historical chat used. Target
median one-ID flow under 20 seconds after project opens; zero unnoticed scope changes.
If the drawer is hard to discover, improve its affordance and pending-state visibility.
Do not substitute the rejected top-only layout. Top session tabs are navigation and do
not imply bringing back a top-only Context Bag editor.


## Parameterized dashboards: one input, several answers

An **Order overview** dashboard declares required scalar `orderId` with semantic
Order.ID and identity-domain metadata. Opening it in a session with one compatible
applied value prefills the input and shows **Order 123 · Context r4**. A cohort never
silently becomes its first order. Offer choose one order or a separately declared
cohort dashboard. Missing/ambiguous inputs show Needs input; never run all orders.

Reuse existing Board.Parameters and widget query parameter mappings. Board definitions
remain reusable: ordinary input IDs, saved query references, typed constants or board
parameter references. No bag IDs in definitions. Existing requiredParams nested-array
semantics need specification before use; the initial dashboard uses one plainly required
input and does not invent AND/OR interpretation for that legacy field.

Recommended first dashboard has an order-summary table, an OrderDetails grid and two
charts: **order value by product category** and **quantity by product**. Both derive from
order-line data. Use one currency, show uncategorized products, aggregate duplicate lines,
and state whether value is quantity × unit price before tax/discount. For the initial
fixture use gross line value before tax/discount; returns/negative values use a diverging
bar rather than a misleading donut. Every chart has a tabular/text alternative.
The review prototype additionally illustrates a progression chart using explicitly
synthetic transition events, giving three charts total. Product delivery requires at
least one chart; two are recommended, within the owner's 1–5 range.

Open resolves inputs; **Run dashboard** intentionally executes. Resolve/pin inputs once
for the entire run: session, applied manifest, board revision, query revisions, mappings,
environment and time anchor. Each widget maps only declared inputs and receives a
separate execution receipt. Missing optional widgets can show partial state; missing
required inputs block dependent widgets. Shared inputs are not an atomic multi-source
snapshot: display source execution times and partial failures. Charts never turn missing
rows or failed queries into zero-valued success.

After Apply, an open dashboard shows **Context changed · Refresh dashboard**. It keeps
its old pinned run until explicit refresh. Explorer can refresh automatically; dashboards
and Checks have explicit cost/execution control. An opt-in refresh-on-Apply dashboard
policy is a later alternative. Saving a dashboard run to chat or Case keeps the original
widget receipts and input manifest. Moving from dashboard to query to Explorer retains
the same session context.

Future progression chart requires actual transition history (order key, event identity,
state, effective/recorded timestamps). Repeated states, reversals, ties and open stages
need documented duration rules. Current status alone cannot establish past durations;
show History unavailable instead of fabricating it.


## Approved: multiple bags, scoped to environments

The owner approved both session collections and multiple bags in a session, with context
applying only to its associated environment. A session is the investigation; a bag is
one working set. Start with one unnamed/default bag and the verified current environment.
Expose naming and a bag switcher when the second bag is created. **UAT affected orders**,
**PROD affected orders**, and **PROD healthy orders** can belong to the same session.
If environment is ambiguous, choose it explicitly before Apply; never infer from a
physical connection shared by QA and UAT. The initial design requires an associated
environment for automatic application. Unassigned imported bags stay Needs environment.

The selected bag is the one being edited; applied state is a map of environment → pinned
bag manifest. One bag per environment is applied by default. Switching the editor from
PROD affected to PROD healthy stages a choice; Apply replaces only PROD's active slot.
Switching to UAT and applying affects UAT only, so simultaneous UAT and PROD tabs can
retain their respective context. The strip always shows bag and environment plus
Applied/Pending; it never just says “300 IDs” without scope. The master checkbox is explicitly labeled **All session context** and stages global
suspension across every environment until Apply, retaining every slot and entry selection.
Apply also updates the selected bag’s own environment slot; no other environment’s
entries change. A disabled individual bag/entry is distinct from this global gate.

**Compare bags** is explicit execution, not an implicit Apply to the entire workspace.
Choose left/right bags, query or parameterized dashboard and a common time anchor if
needed. Both sides get independently pinned manifests, their own environment/source
mapping, authorization and timestamps. Different environments may be compared side by
side, never unioned by numeric ID. Joining identities across environments requires an
explicit business-key alignment and comparable metric units/currencies. Without one,
show separate labeled charts; do not invent per-order deltas. Same-environment cohort
union/overlap remains supported. PROD/UAT authorization denial must not leak the denied
side’s counts, titles or cached data through the permitted side.

**Clear context → Confirm clear** is exactly two clicks. The confirmation names the
selected bag and environment and explains preservation of historical evidence. Confirm
atomically removes its current entries and unsent bag attachments and invalidates its
active slot only if that slot points to this bag. It leaves another applied bag for the
same environment untouched. The cleared bag keeps its identity/name/environment; its
new empty bag revision is saved and its active slot is removed if it was active. No third Apply is needed.
Other bags, composer text, sent turns, past dashboards and Case evidence remain intact.
Use Delete bag separately to remove the bag itself with reference-aware tombstoning.
Clearing is enabled for running sessions; stopped users Resume first. Button confirmation
itself does not authorize fresh unconstrained queries; affected views show Context cleared
and an explicit View all route, while historical runs remain pinned and labeled.

Early feature set: New bag, Rename, Duplicate snapshot, switch/edit, Apply, Clear and
Compare. Cross-session reusable bag library, shared linked bags and live updates remain
later capabilities. Duplicate/copy creates a new entry identity with origin references;
changing environment cannot simply reparent entries and reinterpret keys.

## Integration with refreshed Query Author

Apps main a35c58e9 (pulled 2026-10-13) already has Parameters / Sources / Related in the
query sidebar, stacked editor/results, mobile Jump to results / Back to editor, row
“Use as filter” with Undo, and execution-details accordions. Place the session strip
above that workspace; the bag drawer temporarily uses inspector space without adding a
permanent fourth column. Preserve code editor width and mobile focus restoration.

Keep **Add to Context Bag** distinct from **Use as filter**: the former stages a semantic
entry, the latter edits query text. Keep query Saved/Unsaved separate from bag
Applied/Pending and store Saved locally/Published/Pending states. Reuse the Parameters
input-recovery UI and execution-details provenance instead of duplicating them. The
standalone prototype illustrates session interactions, not a pixel-faithful replacement
for the refreshed Query Author.

## Component strategy: Ionic first, AG Grid for tabular data

Owner preference: use Ionic wherever applicable; PrimeNG requires a concrete advantage
or a separately evaluated broader shift toward denser content components. The inspected
apps package at a35c58e9 declares Ionic 9.0.0 and AG Grid 36.2.0; its package manifest and
lockfile do not declare PrimeNG. Existing board widgets already use `ion-segment`, and
the semantic context panel uses `ion-item`. Prototype styling is an interaction reference,
not a requirement to reproduce bespoke controls or install a new component library.

**Use `ion-item` as the main row composition primitive wherever applicable**, including
row-like card headers as explicitly requested by the owner. Compose title/supporting
text with `ion-label`, icons at start, and concise status/actions at end. For a card with
a single header row, use a semantic heading inside that row; avoid doubling `ion-card-header`
padding around it. A header with independent action buttons remains non-clickable.
For selectable entry rows with checkbox and Attach/menu actions, use a non-button item;
each control has its own accessible label. Avoid nesting interactive elements inside
an item configured as a button. Use ordinary layout containers for charts/editor canvases
and multi-column workspaces, where a list-row primitive adds no meaning.

| Surface | Preferred building blocks | Application-specific work |
|---|---|---|
| Session tabs | `ion-tabs` with top `ion-tab-bar` and `ion-tab-button` for each open session, as requested by the owner | Session service owns open/focus/close model, draft guard and URL restoration; integrate existing router and verify dynamic-tab overflow |
| Context summary / card headers | `ion-item`, `ion-label`, `ion-note`/`ion-badge`, slot icons and `ion-button` actions | Compact pending/applied state; semantic headings; no clickable header containing buttons |
| Bag picker / All sessions | `ion-popover` or `ion-modal`, `ion-searchbar`, `ion-list` + `ion-item` | Scope/state labels, search, empty state and reopening |
| Entry selection and editing | `ion-item`, `ion-checkbox`, `ion-input`, `ion-select`, `ion-textarea` | Typed editing, validation, provenance and staged selection |
| Composer references | `ion-chip` with labeled remove action | Revision/source disclosure and send pinning |
| Bag drawer | Existing inspector host + Ionic controls; `ion-modal` on narrow screens | Desktop panel placement, focus return, sticky Apply and two-click Clear |
| Large member/result views | Existing AG Grid | Virtualized rows, selection, typed values and scope-specific receipts |
| Confirmations / overflow | `ion-alert`, `ion-popover`, `ion-action-sheet` where appropriate | Clear scope/impact wording and explicit lifecycle commands |
| Dashboard charts | Dedicated chart renderer within Ionic cards | Chart choice evaluated independently; not a reason alone to adopt PrimeNG |

The owner explicitly selected **Ionic tabs in the real application**. Use `ion-tabs`,
a top-positioned `ion-tab-bar`, and `ion-tab-button` session navigation; no custom tab
strip or `ion-segment` substitute for this surface. The plain HTML prototype remains
an interaction illustration and does not need Ionic dependencies. Adapt dynamic session
IDs to the existing Ionic/Angular route structure; the session service remains the
authority for drafts, lifecycle, applied context and retained results. Verify cold links,
Back behavior, content restoration and tab overflow in the actual app instead of assuming
Ionic navigation alone persists the investigation. Close actions remain outside the
interactive tab button; a selected-session actions menu is acceptable if a per-tab close
affordance conflicts with Ionic markup. See [Ionic Tabs](https://ionicframework.com/docs/api/tabs)
and [Ionic Tab Bar](https://ionicframework.com/docs/api/tab-bar).

Achieve desktop density with scoped spacing/height variables, restrained header padding,
one-line primary rows, and expandable secondary metadata. Keep touch controls generous,
full labels accessible, keyboard focus visible, and text usable at 200% zoom. Prefer
documented variables/parts to shadow-DOM overrides. Ionic documents row content/actions
and warns about nested interactive controls in [ion-item](https://ionicframework.com/docs/api/item);
its general recommendation is lists, while the owner's selected card-header composition
reuses the same row pattern. Validate heading and screen-reader behavior in that placement.

PrimeNG is a **conditional alternative**, not part of the Context Bag baseline. Its
[Tabs](https://primeng.dev/tabs) and [theme tokens](https://primeng.dev/theming/styled)
offer useful building blocks, but neither decides our session lifecycle, density or
accessibility for us. Before adding it, demonstrate a concrete Ionic limitation on the
same screen and compare keyboard/touch behavior, 200% zoom, dark mode, incremental
production bundle size, Angular compatibility and maintenance cost. Do not invent a
bundle penalty without measuring it. Keep AG Grid instead of adding a second grid.

If several major content surfaces need replacements, bring back a separate decision
for Ionic navigation/layout plus mostly PrimeNG content controls, with a shared density,
theme, overlay/focus and migration policy. Avoid gradually accumulating two versions
of buttons, inputs, menus and dialogs without that decision. First validation should
use the session strip, bag header/list and drawer in actual Ionic after implementation
approval; the standalone HTML prototype does not prove framework fit.
