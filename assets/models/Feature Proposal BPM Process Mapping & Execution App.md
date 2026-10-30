# Feature Proposal: BPM Process Mapping & Execution App

Proposed

## 1. Summary

This proposal describes an application for **mapping, running and tracking business processes**. Processes are described in a simple, versionable text format (Markdown + YAML), rendered automatically as BPM flows, and can be drilled into at the activity level: each activity can collect data, require validation, call a backend API, or trigger an ETL workflow. Activities can be performed manually by assigned actors or run autonomously on a schedule.

The application has two layers:

1. **Workflow definitions**: the process model (nodes, edges, actors, data schemas, actions).
2. **Workflow execution**: running instances of a process, step tracking, and all data collected along the way.

## 2. Goals

- Let an organization model its business processes, including macro processes and nested sub-processes.
- Generate the BPM diagram automatically from the text definition.
- Go beyond the diagram: attach data collection, validation and automation to each activity.
- Support manual execution by stakeholders, automated execution, or a mix of both.
- Notify the right person when it is their turn, and keep a full audit trail of every step and every piece of data collected.

## 3. Core Concepts

In an organization there are many **Business Processes (BP)**. A BP can live inside a macro BP or contain sub-processes, which are themselves BPs. Every process begins at a **Starting Point (SP)**. A BP is made of **activities (BA)**, the steps through which work gets done. Activities are performed by **stakeholders (ST)** (actors), each with a lane (horizontal or vertical). Flows may branch at **Decision Points (DP)**, which route to another activity, another process, or an end.

### 3.1 Element codes

| Prefix | Element | Code logic |
| --- | --- | --- |
| `BP` | Business Process | Auto-incremented number, e.g. `BP001` |
| `BA` | Business Activity / Task | BP sequence + activity sequence inside that BP |
| `SP` | Starting Point | Same logic as `BA` |
| `DP` | Decision Point | Same logic as `BA`; takes three arguments: `(condition, flow_if_true, flow_if_false)` |
| `EP` | End Point | Same logic as `BA` |
| `ST` | Stakeholder / Actor | Same logic as `BA` |
| `FL` | Flow | Connects the points above |

## 4. Definition Format

Processes are written in Markdown, with one YAML block per element. Heading levels carry the structure (`#` model, `##` stakeholders and processes, `##` activities).

````markdown
 #BPM
```yaml
name: BPM
runs_as: BPM
```
 
##ST001
```yaml
name: DEI
desc: Department EI
focal_points:
  - {name: x, email: userx@domain.com, role: x} 
  - {name: y, email: usery@domain.com, role: y}
```
 
##ST002
```yaml
name: DEE
desc: Department EE
focal_points:
  - {name: z, email: userz@domain.com, role: z} 
```
 
## BP001
```yaml
name: Process 1
desc: Desc Process 1
```
 
## BA001
```yaml
name: Act 1
desc: Act 1 of Process 1
actors:
   - ST001
data_schema: 
  - {name: field_x, type: varchar, desc: shdhgs, ..., size: 10}
  - {name: field_y, type: bool, desc: shdhgs, ..., size: 2}
```
 
## BA002
```yaml
name: Act 2
desc: Act 2 of Process 1
actors:
   - ST001
data_schema: 
  - {name: field_x, type: varchar, desc: shdhgs, ..., size: 10}
  - {name: field_y, type: bool, desc: shdhgs, ..., size: 2}
```
 
## BA003
```yaml
name: Act 3
desc: Act 3 of Process 1
actors:
   - ST002
data_schema: 
  - {name: field_x, type: varchar, desc: shdhgs, ..., size: 10}  
  - {name: attach_x, type: varchar, desc: attachments, ..., size: 10}
  - {name: field_y, type: bool, desc: shdhgs, ..., size: 2}
action_buttons: # call api with the data
  - {name, label, type, api: etlx/run/daily-report-y}
```
 
## FL001
```yaml
   - SP001->BA001->BA002->DP001(BA002.field_y=TRUE,BA003,BA001)->EP001 
```
````

In this example, `DP001` evaluates `BA002.field_y = TRUE`: if true the flow continues to `BA003`, otherwise it returns to `BA001`.

## 5. Data Model: Workflow Definitions

The definition layer works like a graph database: a **nodes** table and an **edges** table, plus supporting tables.

### 5.1 `process_node`

Holds every element of the model, equivalent to the nodes table of a graph database.

- `id`, `code`, `description`, `attachments` (documents related to the object)
- The code prefix determines the node type (`BP`, `BA`, `SP`, `DP`, `EP`, `ST`).

### 5.2 Stakeholders and focal points

For every edge touching an `ST***` node, the system looks up that stakeholder's `focal_points`. When the SP generates its id and code, the focal points are populated with their name and contact details.

If a focal point is also a user of the app, **their email must match their user email**. This is what allows the system to notify them when the previous step ends and the next one belongs to them, so they can pick up where the previous actor left off.

### 5.3 Activity details

For every `BA***` registered in the nodes table, the following are stored in their own tables:

- **Actors**: which stakeholders perform the activity.
- **Data schema**: the fields to collect, with types, sizes and validation rules.
- **Action buttons**: API calls or ETL runs that use the collected data.

### 5.4 Flows (process edges)

Flows are not stored directly as edges. They are defined at node level (`FL***`) and then **expanded into an edges table** (called `process_edge` or `flows`), exactly as edges work in a graph database.

The parsing rule:

1. Read the flow as an array of strings.
2. Split each string on `->`.
3. Register each consecutive pair (`node1`, `node2`) as an edge.

### 5.5 Decision points

Decision points need special handling. When a `DP` is found, its two outgoing paths are registered as **conditional edges**:

- Edge to the *true* target, flagged `is_decision = true` with the condition and the branch value (true).
- Edge to the *false* target, flagged the same way with the branch value (false).

Both paths exist in a kind of superposition: at runtime the condition is evaluated and only the matching path is activated and displayed.

## 6. Data Model: Workflow Instantiation and Execution

Only **processes** are shown to end users. When a user clicks **Start**, the system creates an **instance** of that process and tracks every step.

### 6.1 Instance tracking

A table records, for each instance:

- which process it belongs to,
- who started it and when,
- the current step/activity and the actor currently responsible.

Every time someone completes an activity, step or sub-process, the instance advances to the next step, and the next actor is notified.

### 6.2 Set-backs

Steps can be sent back:

- manually by the **process coordinator**, or
- automatically by a **decision point** whose condition routes the flow backwards.

### 6.3 Collected data

All data collected at every step is saved with a reference to the **instance**, **process**, **activity** and **actor**, giving a complete audit trail.

## 7. Execution Modes

- **Manual**: an actor opens their task, fills in the data and completes it.
- **Automated**: the activity calls a backend API or runs an ETL workflow, either on a schedule or when it is reached.
- **Mixed**: a person provides data and then triggers an action button.

## 8. Open Questions

- **Code format**: the description says `BA******` (BP sequence + activity sequence, e.g. `BA001001`), but the examples use `BA001`. Which convention should we adopt?
- **Condition language**: what syntax is allowed in decision conditions (`BA002.field_y=TRUE`), and how is it evaluated safely?
- **Lanes**: how are horizontal/vertical lane preferences stored: on the stakeholder or on the diagram?
- **Sub-process completion**: does the parent activity complete automatically when the child process reaches its end point?
- **Scheduling**: where do the schedules for autonomous activities live, and who can edit them?
- **Versioning**: what happens to running instances when a process definition changes?
