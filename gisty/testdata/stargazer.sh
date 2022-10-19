#!/bin/bash
# shellcheck disable=SC2016

# Using the `gh api graphql` command to retrieve and format information
QUERY1='
query {
  viewer {
    gists (first: 10, orderBy: {field: CREATED_AT, direction: DESC} ) {
        nodes {
            createdAt
            description
            name
            pushedAt
            stargazers (first: 100) {
            totalCount
            edges {
                node {
                id
                }
            }
            }
            updatedAt
        }
    }
  }
}
'
TEMPLATE1='
  {{- range $repo := .data.viewer.gists.nodes -}}
    {{- printf "name: %s - stargazers: %v\n" $repo.name $repo.stargazers.totalCount -}}
  {{- end -}}
'

# QUERY2='
# query {
#   viewer {
#     gist (name: "5b10b34f87955dfc86d310cd623a61d1" ) {
#         name
#         stargazerCount
#     }
#   }
# }
# '

QUERY2='query { viewer { gist (name: "5b10b34f87955dfc86d310cd623a61d1" ) { name, stargazerCount } } }'
# TEMPLATE2='
#     {{- printf "name: %s - stargazers: %v\n" .data.viewer.gist.name .data.viewer.gist.stargazerCount -}}
# '
TEMPLATE2='{{.data.viewer.gist.stargazerCount}}'

gh api graphql -f query="${QUERY1}" --paginate --template="${TEMPLATE1}"
echo "----------------------------------"
gh api graphql -f query="${QUERY2}" --paginate --template="${TEMPLATE2}"