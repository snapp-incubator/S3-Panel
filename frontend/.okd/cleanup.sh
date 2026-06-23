#bin/bash

exclude=/tmp/exclude.txt

echo "getting branches ..."

curl --header "PRIVATE-TOKEN: ${GITLAB_TOKEN}" --silent --fail "https://gitlab.snapp.ir/api/v4/projects/${CI_PROJECT_ID}/merge_requests?state=opened&per_page=100" | python -mjson.tool | grep -w source_branch | tr '[:upper:]' '[:lower:]' | awk -F'"' '{print $4}' | sed -e 's/\//-/g' > ${exclude}

echo "excluding live-tests deployments"
echo "live-test-prod" >> $exclude;
echo "live-test-develop" >> $exclude;
echo "qa-automation-test" >> $exclude;

echo "login to okd ..."

# oc login https://okd.private.teh-1.snappcloud.io --token=${OKD_TEH1_TOKEN}

oc project ${PROJECT_STAGING}

tags=$(oc get ImageStream ${CI_PROJECT_NAME} -o yaml | grep -w tag | grep -v "tag: v." | grep -v develop | grep -v -f ${exclude} | sed "s/tag: /${CI_PROJECT_NAME}:/g")

if [ -z "$tags" ]; then
    echo "No unused tags to clean from image stream!";
    exit 0;
fi

printf "%s\n" "${tags[@]}" | xargs -L1 oc tag -d

echo "delete old deployments"

oc get dc | grep -v develop | grep -v -f ${exclude} | awk 'NR>1{print $1}' | xargs -L1 oc delete dc

echo "delete old configmaps"

oc get cm | grep -v develop | grep -v -f ${exclude} | awk 'NR>1{print$1}' | xargs -L1 oc delete cm

echo "delete old routes"

oc get route | grep -v develop | grep -v -f ${exclude} | awk 'NR>1{print$1}' | xargs -L1 oc delete route

echo "delete old services"

oc get service | grep -v develop | grep -v -f ${exclude} | awk 'NR>1{print$1}' | xargs -L1 oc delete service
