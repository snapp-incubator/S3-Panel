#!/bin/bash

function generate_hash() {
  sha256sum $1 | awk {'print $1'}
}

function env_to_object() {
  local default_delimiter="="

  local target="$1"
  local delimiter="${2:-$default_delimiter}"

  local content=""

  while IFS= read -r line; do
    local key=$(echo $line | awk -F"$delimiter" '{print $1}')
    local value=$(echo $line | awk -F"$delimiter" '{print $2}')
    local entry="$key:\"$value\","
    content=$content$entry
  done <$target

  echo "{$content}"
}

function create_config_map() {
  local target_directory="$1"

  local variables=$(ls ${target_directory})
  local config_map=""

  for key in $variables; do
    value=$(cat ${target_directory}/${key})
    config_map+="${key}=${value}"
    config_map+=$'\n'
  done

  echo "$config_map"
}

function remove_blank_lines() {
  local target="$1"
  sed -i '/^$/d' $target
}

# Create config-map from the files inside /etc/runtime-variables directory
create_config_map "/etc/runtime-variables" >>/tmp/.config-map

# Remove all blank lines before generating an object
remove_blank_lines /tmp/.config-map

# Add public configs into the config-map as well
printenv | grep ^PUBLIC_CONFIG_ >>/tmp/.config-map

# Replace configMap placeholders with actual values.
sed -i -e "s|__SERVICE_NAMESPACE__|$SERVICE_NAMESPACE|g" -e "s|__CLUSTER_NAME__|$CLUSTER_NAME|g" /tmp/.config-map

# Convert text based configMap to object.
config_map=$(env_to_object /tmp/.config-map)

# Insert configMap into the HTML's config script.
sed -i -e 's|id="config">|id="config">window.configuration='"$config_map"'|' /usr/share/nginx/html/index.html

nginx -g 'daemon off;'
