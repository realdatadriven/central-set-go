# SQL Schema for Custom Survey Builder

This document outlines a relational SQL schema designed to support a custom survey builder application. The schema allows for the creation of surveys with multiple pages, various question types, and conditional logic, enabling a flexible and robust survey design experience.

## 1. `Surveys` Table
This table stores general information about each survey.

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `survey_id` | `INT` | `PRIMARY KEY`, `AUTO_INCREMENT` | Unique identifier for the survey. |
| `title` | `VARCHAR(255)` | `NOT NULL` | The title of the survey. |
| `description` | `TEXT` | `NULLABLE` | A brief description of the survey. |
| `created_at` | `DATETIME` | `DEFAULT CURRENT_TIMESTAMP` | Timestamp when the survey was created. |
| `updated_at` | `DATETIME` | `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | Timestamp when the survey was last updated. |
| `status` | `VARCHAR(50)` | `DEFAULT 'draft'` | Current status of the survey (e.g., 'draft', 'published', 'archived'). |

## 2. `Pages` Table
This table defines the structure of pages within a survey, allowing for multi-page surveys.

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `page_id` | `INT` | `PRIMARY KEY`, `AUTO_INCREMENT` | Unique identifier for the page. |
| `survey_id` | `INT` | `NOT NULL`, `FOREIGN KEY (Surveys.survey_id)` | Foreign key linking to the `Surveys` table. |
| `page_order` | `INT` | `NOT NULL` | The order of the page within the survey. |
| `title` | `VARCHAR(255)` | `NULLABLE` | The title of the page. |
| `description` | `TEXT` | `NULLABLE` | A brief description or instruction for the page. |

## 3. `Questions` Table
This table stores individual questions, linking them to specific pages.

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `question_id` | `INT` | `PRIMARY KEY`, `AUTO_INCREMENT` | Unique identifier for the question. |
| `page_id` | `INT` | `NOT NULL`, `FOREIGN KEY (Pages.page_id)` | Foreign key linking to the `Pages` table. |
| `question_order` | `INT` | `NOT NULL` | The order of the question within its page. |
| `text` | `TEXT` | `NOT NULL` | The actual question text. |
| `question_type_id` | `INT` | `NOT NULL`, `FOREIGN KEY (QuestionTypes.question_type_id)` | Foreign key linking to the `QuestionTypes` table. |
| `is_required` | `BOOLEAN` | `DEFAULT FALSE` | Indicates if the question is mandatory. |
| `placeholder` | `VARCHAR(255)` | `NULLABLE` | Placeholder text for input fields. |
| `default_value` | `TEXT` | `NULLABLE` | Default value for the question. |

## 4. `QuestionTypes` Table
This is a lookup table for predefined question input types.

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `question_type_id` | `INT` | `PRIMARY KEY`, `AUTO_INCREMENT` | Unique identifier for the question type. |
| `type_name` | `VARCHAR(50)` | `NOT NULL`, `UNIQUE` | Name of the question type (e.g., 'text', 'textarea', 'number', 'radio', 'checkbox', 'dropdown', 'date', 'email'). |

## 5. `Choices` Table
This table stores predefined options for questions like radio buttons, checkboxes, and dropdowns.

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `choice_id` | `INT` | `PRIMARY KEY`, `AUTO_INCREMENT` | Unique identifier for the choice. |
| `question_id` | `INT` | `NOT NULL`, `FOREIGN KEY (Questions.question_id)` | Foreign key linking to the `Questions` table. |
| `value` | `VARCHAR(255)` | `NOT NULL` | The internal value of the choice. |
| `text` | `VARCHAR(255)` | `NOT NULL` | The display text for the choice. |
| `choice_order` | `INT` | `NOT NULL` | The order of the choice within the question. |

## 6. `Conditions` Table
This table defines conditional logic for questions or pages. For simplicity, this initial design focuses on showing/hiding elements based on a single condition. More complex logic (AND/OR groups) would require additional tables or a more sophisticated JSON structure within a `condition_json` field.

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `condition_id` | `INT` | `PRIMARY KEY`, `AUTO_INCREMENT` | Unique identifier for the condition. |
| `target_type` | `VARCHAR(50)` | `NOT NULL` | 'question' or 'page' to indicate what the condition applies to. |
| `target_id` | `INT` | `NOT NULL` | ID of the question or page being conditionally displayed. |
| `source_question_id` | `INT` | `NOT NULL`, `FOREIGN KEY (Questions.question_id)` | The question whose answer triggers the condition. |
| `operator` | `VARCHAR(50)` | `NOT NULL` | Comparison operator (e.g., '=', '!=', '>', '<', 'contains'). |
| `value` | `TEXT` | `NOT NULL` | The value to compare against the source question's answer. |
| `action` | `VARCHAR(50)` | `NOT NULL` | 'show' or 'hide' the target element. |

## Relationships

*   `Surveys` 1 -- M `Pages`
*   `Pages` 1 -- M `Questions`
*   `Questions` M -- 1 `QuestionTypes`
*   `Questions` 1 -- M `Choices` (for choice-based questions)
*   `Questions` 1 -- M `Conditions` (as `source_question_id`)
*   `Pages` 1 -- M `Conditions` (as `target_id` when `target_type` is 'page')
*   `Questions` 1 -- M `Conditions` (as `target_id` when `target_type` is 'question')

This schema provides a foundation for building a survey creator. The next steps will involve defining the specific input types and detailing how conditional logic can be structured within this model.

## 7. Input Types and Their Attributes

The `QuestionTypes` table provides a basic categorization. However, different input types will have specific attributes that need to be considered when rendering the form. These attributes can be stored as a JSON string in a `settings` column within the `Questions` table, or in a separate `QuestionSettings` table for more structured storage.

Here are common input types and their typical attributes:

| `type_name` (from `QuestionTypes`) | Description | Common Attributes |
| :-------------------------------- | :---------- | :---------------- |
| `text` | Single-line text input. | `min_length`, `max_length`, `regex_pattern` |
| `textarea` | Multi-line text input. | `min_length`, `max_length`, `rows`, `cols` |
| `number` | Numeric input. | `min_value`, `max_value`, `step` |
| `radio` | Single choice from a list. | `choices` (from `Choices` table) |
| `checkbox` | Multiple choices from a list. | `choices` (from `Choices` table), `min_selected`, `max_selected` |
| `dropdown` | Single choice from a dropdown list. | `choices` (from `Choices` table), `is_searchable` |
| `date` | Date picker. | `min_date`, `max_date`, `date_format` |
| `email` | Email input with validation. | `regex_pattern` (for custom validation) |
| `url` | URL input with validation. | `regex_pattern` |
| `file` | File upload. | `allowed_file_types`, `max_file_size`, `max_files` |
| `rating` | Star or numerical rating. | `min_rating`, `max_rating`, `step` |
| `slider` | Range slider. | `min_value`, `max_value`, `step` |

For example, the `Questions` table could be extended with a `settings_json` column:

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `settings_json` | `JSON` or `TEXT` | `NULLABLE` | JSON string containing type-specific settings (e.g., `{"min_length": 5, "max_length": 100}`). |

## 8. Conditional Logic Structure

The `Conditions` table, as defined, supports basic show/hide logic based on a single source question's answer. To implement more complex conditional logic, such as combining multiple conditions with AND/OR operators, or more advanced actions (e.g., setting a value, making a question required), the `Conditions` table can be enhanced.

### Enhanced `Conditions` Table (Conceptual)

| Column Name | Data Type | Constraints | Description |
| :---------- | :-------- | :---------- | :---------- |
| `condition_id` | `INT` | `PRIMARY KEY`, `AUTO_INCREMENT` | Unique identifier for the condition. |
| `target_type` | `VARCHAR(50)` | `NOT NULL` | 'question' or 'page'. |
| `target_id` | `INT` | `NOT NULL` | ID of the question or page being conditionally affected. |
| `logic_json` | `JSON` or `TEXT` | `NOT NULL` | JSON string defining the complex conditional logic. |
| `action` | `VARCHAR(50)` | `NOT NULL` | 'show', 'hide', 'enable', 'disable', 'set_value', 'set_required'. |
| `action_value` | `TEXT` | `NULLABLE` | Value to set if action is 'set_value'. |

### Example `logic_json` Structure for Complex Conditions

For a condition like "Show Question C if (Question A = 'Yes' AND Question B > 10) OR (Question A = 'Maybe')", the `logic_json` could look like this:

```json
{
  "operator": "OR",
  "rules": [
    {
      "operator": "AND",
      "rules": [
        {
          "source_question_id": 1, 
          "comparison_operator": "=",
          "value": "Yes"
        },
        {
          "source_question_id": 2,
          "comparison_operator": ">",
          "value": 10
        }
      ]
    },
    {
      "source_question_id": 1,
      "comparison_operator": "=",
      "value": "Maybe"
    }
  ]
}
```

This `logic_json` field would allow for a highly flexible and nested conditional logic system, which your application would parse and apply when rendering the survey. The `source_question_id` would reference questions in the `Questions` table, and `comparison_operator` and `value` would define the specific condition to check against the answer of the `source_question_id`.

This structured approach ensures that all necessary information for rendering and controlling survey flow is stored relationally, yet allows for the flexibility of complex configurations through JSON fields.

## 9. Merging SQL Data into a JSON Model

The goal is to transform the relational data stored in the SQL tables into a single, hierarchical JSON object that can be easily consumed by a front-end application to render the survey. This process typically involves a series of SQL queries to fetch the data and then programmatic assembly of the JSON structure in your application's backend (e.g., using Python, Node.js, PHP).

### Merging Strategy

1.  **Fetch Survey Details**: Start by querying the `Surveys` table for the main survey information.
2.  **Fetch Pages**: For the given `survey_id`, fetch all associated pages from the `Pages` table, ordered by `page_order`.
3.  **Fetch Questions per Page**: For each page, fetch its associated questions from the `Questions` table, ordered by `question_order`.
4.  **Fetch Question Details**: For each question:
    *   Retrieve its `question_type_id` and join with `QuestionTypes` to get `type_name`.
    *   Parse the `settings_json` (if present) into an object.
    *   If the question type is `radio`, `checkbox`, or `dropdown`, fetch its `Choices` from the `Choices` table, ordered by `choice_order`.
5.  **Fetch Conditions**: Fetch all relevant conditions from the `Conditions` table. These will need to be processed to form the `logic_json` structure if complex conditions are used, or simpler `visibleIf` properties for basic conditions.

### Example JSON Output Structure

This example demonstrates how the SQL data could be transformed into a JSON structure, inspired by common survey JSON formats. Note that specific property names (e.g., `elements`, `choices`, `visibleIf`) can be adapted to fit your front-end framework's requirements.

```json
{
  "survey_id": 1,
  "title": "Customer Feedback Survey",
  "description": "We appreciate your feedback!",
  "status": "published",
  "pages": [
    {
      "page_id": 101,
      "title": "About You",
      "description": "Tell us a little about yourself.",
      "elements": [
        {
          "question_id": 1001,
          "type": "text",
          "name": "fullName",
          "text": "What is your full name?",
          "is_required": true,
          "placeholder": "John Doe",
          "settings": {
            "max_length": 100
          }
        },
        {
          "question_id": 1002,
          "type": "email",
          "name": "emailAddress",
          "text": "What is your email address?",
          "is_required": false
        },
        {
          "question_id": 1003,
          "type": "radio",
          "name": "gender",
          "text": "What is your gender?",
          "is_required": true,
          "choices": [
            { "value": "male", "text": "Male" },
            { "value": "female", "text": "Female" },
            { "value": "other", "text": "Other" }
          ]
        }
      ]
    },
    {
      "page_id": 102,
      "title": "Your Experience",
      "description": "Share your thoughts on our service.",
      "elements": [
        {
          "question_id": 2001,
          "type": "rating",
          "name": "overallRating",
          "text": "How would you rate your overall experience?",
          "is_required": true,
          "settings": {
            "min_rating": 1,
            "max_rating": 5
          }
        },
        {
          "question_id": 2002,
          "type": "textarea",
          "name": "comments",
          "text": "Any additional comments?",
          "is_required": false,
          "visibleIf": "{overallRating} < 3",
          "settings": {
            "rows": 4
          }
        },
        {
          "question_id": 2003,
          "type": "checkbox",
          "name": "featuresUsed",
          "text": "Which features did you use?",
          "is_required": false,
          "choices": [
            { "value": "featureA", "text": "Feature A" },
            { "value": "featureB", "text": "Feature B" },
            { "value": "featureC", "text": "Feature C" }
          ]
        }
      ]
    }
  ],
  "conditions": [
    {
      "target_type": "question",
      "target_id": 2002,
      "logic": {
        "source_question_id": 2001,
        "comparison_operator": "<",
        "value": 3
      },
      "action": "show"
    }
  ]
}
```

This JSON structure provides a comprehensive representation of the survey, including its pages, questions, choices, and conditional logic, all derived from the proposed SQL schema. Your application would then interpret this JSON to dynamically render the survey form to the user.
