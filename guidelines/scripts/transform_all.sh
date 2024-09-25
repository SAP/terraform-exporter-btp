#!/bin/bash

# Create the plan file
terraform plan -out=plan.out

# Generate the JSON file
terraform show -json plan.out |  jq .planned_values.outputs.all.value > restrictedplan.json

# Read the JSON file
json_content=$(cat restrictedplan.json)

# Mask every " with \"
masked_content=$(echo "$json_content" | sed 's/"/\\"/g')

# Remove all line breaks
no_line_breaks=$(echo "$masked_content" | tr -d '\n')

# Remove spaces before or after a masked quote
no_spaces_around_quotes=$(echo "$no_line_breaks" | sed 's/ *\\\" */\\"/g')

# Remove spaces between any kind of brackets, including nested ones
transformed_content=$(echo "$no_spaces_around_quotes" | sed -E 's/\{[[:space:]]+/\{/g; s/[[:space:]]+\}/\}/g; s/\[[[:space:]]+/\[/g; s/[[:space:]]+\]/\]/g; s/\([[:space:]]+/\(/g; s/[[:space:]]+\)/\)/g')

# Write the transformed content to a new text file
echo "$transformed_content" > transformed_restrictedplan.txt

# Remove the temporary plan file
rm restrictedplan.json
