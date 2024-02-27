# Use Github Page to Configure public helm chart repository

## Create a new GitHub repository

## Clone the repository
```shell
git clone git@github.com:ChrisDevOpsOrg/helm-charts.git
cd helm-charts
```

## Create robots.txt
### To avoid bot crawling
```shell
echo -e "User-Agent: *\nDisallow: /" > robots.txt
```

## Create a helm chart from scratch
```shell
mkdir -p chart-sources
cd chart-sources
helm create demo
```

## lint the chart
```shell
helm lint chart-sources/*
```

## Create helm chart package
```shell
helm package chart-sources/nginx-demo
```
## Configure the “helm-chart” repository as a GitHub pages site
Go to GitHub repository -> settings -> GitHub Pages

## Create the helm chart repository index
> A repository is characterized primarily by the presence of a special file called index.yaml that has a list of all of the packages supplied by the repository, together with metadata that allows retrieving and verifying those packages.
```shell
# The URL is your github page URL, is not github repo URL
helm repo index --url https://ChrisDevOpsOrg.github.io/helm-charts/ .
```

## Push your charts to Github
```shell
git add . && git commit -m "Initial commit" && git push origin master
```

## Add new charts to an existing repository
> Whenever you need to add a new chart to the Helm chart repository, it’s mandatory for you to regenerate the index.yaml file. The $ helm repo index command will completely rebuild the index.yaml file from scratch, including only the charts that it finds locally, which very likely is our case. However, it worth notice that you can use the --merge flag to incrementally add new charts to an existing index.yaml
```shell
helm repo index --url https://ChrisDevOpsOrg.github.io/helm-charts/ --merge index.yaml .
```

---
# Configure Github private Repo as helm chart repository

## Create helm chart repository index
> For the private helm chart repo, you dont need to specify URL.
```shell
helm repo index .
```
Your index.yaml will looks like:
```yaml
apiVersion: v1
entries:
  demo:
  - apiVersion: v2
    appVersion: 1.0.0
    created: "2023-11-25T08:21:50.859702701+08:00"
    description: A Helm chart for Kubernetes
    digest: 0acf4ab914823060cbd07e104105e733a5fc56fbfc15ecd86aece6b9e831c05e
    name: demo
    type: application
    urls:
    - demo-0.1.1.tgz
    version: 0.1.1
  - apiVersion: v2
    appVersion: 1.0.0
    created: "2023-11-25T08:21:50.85941036+08:00"
    description: A Helm chart for Kubernetes
    digest: f88a2a39923ae48e42c314c521254b4288a51f807c6b90e4041d81cc1ed02e25
    name: demo
    type: application
    urls:
    - demo-0.1.0.tgz
    version: 0.1.0
generated: "2023-11-25T08:21:50.8590034+08:00"
```

## Installing from the Private Repository
```shell
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN"
helm repo add demo-charts "https://${GITHUB_TOKEN}@raw.githubusercontent.com/ChrisDevOpsOrg/helm-charts/main/"
```
> - Enclose the url with double quote " "
> - The trailing / is mandatory
> - If you want to use different branch, change the **main** to your branch name.

Onec the repos is added, you can search it or install charts from it. Note that you will have to update the local
repository index when looking for new versions.
```shell
helm repo update
helm search repo demo-charts --devel
```
