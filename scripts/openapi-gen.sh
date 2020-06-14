#! /bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

which yq >/dev/null 2>&1
YQ_EXISTS=$?
if [ $YQ_EXISTS -ne 0 ]; then
  echo "yq is missing..."
  echo "perhaps try 'brew install yq'"
  echo "or install directly from source https://github.com/mikefarah/yq"
  exit -1
fi

which oapi-codegen >/dev/null 2>&1
OAPIGEN_GO_EXISTS=$?
if [ $OAPIGEN_GO_EXISTS -ne 0 ]; then
  echo "oapi-codegen is missing..."
  echo "check https://github.com/deepmap/oapi-codegen for installation directions"
  exit -1
fi

which openapi-generator >/dev/null 2>&1
OAPIGEN_NODE_EXISTS=$?
if [ $OAPIGEN_NODE_EXISTS -ne 0 ]; then
  echo "openapi-generator-cli is missing..."
  echo "check https://github.com/openapitools/openapi-generator-cli for installation directions"
  exit -1
fi

dir_name="$(dirname $0)"
source ${dir_name}/common.sh

GEN_INFO_JSON=$(yq r ${REPO_ROOT}/api/${APP_NAME}.yaml --collect 'x-oapi-codegen.*.generator' -j)
GEN_COUNT=$(echo $GEN_INFO_JSON | yq r - --collect -l '*')

# genOpenApiBoilerplate calls generate with a single generator object
# defined in spec in a 'x-oapi-codegen' field.
function genOpenApiBoilerplate() {
  local generatorObject=$1
  generate "$generatorObject"
}

function generate() {
  local spec=$1
  local language=$(echo $spec | yq r - "language")
  local target=$(echo $spec | yq r - "target")

  case $language in
  "Go" | "go" | "Golang" | "golang")
    local go_package=$(echo $spec | yq r - "go_package")
    local package_path_relative=$(echo $spec | yq r - "path")
    local optins=$(echo $spec | yq r - "opt-ins")

    if [ -z "$go_package" ] || [ -z "$package_path_relative" ]; then
      echo -e "fields 'go_package' and 'path' must be set"
      return -1
    fi
    if [ ! -z $optins ]; then
      # opt-ins are not empty so add ',' at the end of the list
      optins="${optins},"
    fi

    # create new directory if needed
    local package_dir=$(dirname $package_path_relative)
    if [[ ! -d ${REPO_ROOT}/${package_dir} ]]; then
      mkdir -p ${REPO_ROOT}/${package_dir}
    fi

    if [ "$target" == "server" ]; then
      local gen="${optins}server"
      oapi-codegen \
        -package ${go_package} \
        -generate ${gen} \
        -o ${REPO_ROOT}/${package_path_relative} \
        ${REPO_ROOT}/api/${APP_NAME}.yaml

    elif [ "$target" == "client" ]; then
      local gen="${optins}client"
      oapi-codegen \
        -package ${go_package} \
        -generate ${gen} \
        -o ${REPO_ROOT}/${package_path_relative} \
        ${REPO_ROOT}/api/${APP_NAME}.yaml

    else
      echo -e "target: '$target' is unsupported"
      return -1
    fi
    ;;

  "typescript-axios")
    local relative_dir_path=$(echo $spec | yq r - "path")
    # create new directory if needed
    if [[ ! -d ${REPO_ROOT}/${relative_dir_path} ]]; then
      mkdir -p ${REPO_ROOT}/${relative_dir_path}
    fi

    openapi-generator generate \
      --skip-validate-spec \
      -g ${language} \
      -i ${REPO_ROOT}/api/${APP_NAME}.yaml \
      -o ${REPO_ROOT}/${relative_dir_path} >/dev/null
    ;;

  *)
    echo -e "unsupported language"
    return -1
    ;;
  esac
}

for ((i = 0; i < $GEN_COUNT; i++)); do
  generator=$(echo $GEN_INFO_JSON | yq r - "[$i]")
  genOpenApiBoilerplate "$generator"
done
