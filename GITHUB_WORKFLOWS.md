# GitHub Actions Workflow Integration

## Overview

The CI/CD pipelines have been optimized to avoid building the Docker image twice. The Docker image is now built once in the `dockerimage.yml` workflow and reused in the `k8s-deployment.yml` workflow.

## Workflow Structure

### 1. Docker Build and Security Scan (`dockerimage.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Tags starting with `v*`

**Jobs:**
1. **lint** - Code quality and linting checks
2. **test** - Run tests on multiple Go versions (1.24.x, 1.25.x)
3. **security-scan** - Run vulnerability and security scans
4. **build** - Build Docker image and upload as artifact
   - Builds multi-platform image (amd64, arm64) and pushes to registry
   - Builds single-platform image (amd64 only) for testing
   - Exports the test image as a tar file artifact
   - Artifact is available for 1 day
5. **deploy-staging** - Deploy to staging (only on main branch)
6. **deploy-production** - Deploy to production (only on tags)

### 2. Kubernetes Deployment Pipeline (`k8s-deployment.yml`)

**Triggers:**
- Automatically runs after `dockerimage.yml` completes successfully
- Can be triggered manually via `workflow_dispatch`

**Jobs:**
1. **helm-test** - Helm chart linting and unit tests
2. **test-deployment** - Test deployment in Kind cluster
   - Downloads the Docker image artifact from `dockerimage.yml`
   - Loads the image into Kind cluster (no rebuild needed!)
   - Deploys using Helm
   - Tests application endpoints
3. **deploy-staging** - Deploy to staging Kubernetes cluster
4. **deploy-production** - Deploy to production Kubernetes cluster

## Key Changes

### dockerimage.yml Changes:

1. **Separate build steps:**
   ```yaml
   # Multi-platform build for registry
   - Build and push (linux/amd64, linux/arm64)

   # Single-platform build for testing
   - Build for testing (linux/amd64 only)
   - Export to /tmp/image.tar
   ```

2. **Upload artifact:**
   ```yaml
   - name: Upload Docker image artifact
     uses: actions/upload-artifact@v4
     with:
       name: docker-image
       path: /tmp/image.tar
       retention-days: 1
   ```

### k8s-deployment.yml Changes:

1. **New trigger mechanism:**
   ```yaml
   on:
     workflow_run:
       workflows: ["Docker Build and Security Scan"]
       types:
         - completed
       branches: [main]
     workflow_dispatch:
   ```

2. **Download artifact instead of building:**
   ```yaml
   - name: Download Docker image artifact
     uses: dawidd6/action-download-artifact@v3
     with:
       workflow: dockerimage.yml
       name: docker-image
       path: /tmp
   ```

3. **Load pre-built image:**
   ```yaml
   - name: Load Docker image into Kind
     run: |
       docker load -i /tmp/image.tar
       IMAGE_NAME=$(docker images --format "{{.Repository}}:{{.Tag}}" | head -n 1)
       docker tag $IMAGE_NAME ${{ env.APP_NAME }}:test
       kind load docker-image ${{ env.APP_NAME }}:test --name ${{ env.CLUSTER_NAME }}
   ```

## Benefits

✅ **No duplicate builds** - Image is built once and reused
✅ **Faster execution** - Kubernetes deployment starts with pre-built image
✅ **Consistency** - Same image tested and deployed everywhere
✅ **Cost savings** - Reduced CI/CD minutes usage
✅ **Better caching** - Build cache is utilized efficiently

## Workflow Execution Flow

```
Push to main/develop or create tag
    ↓
Docker Build and Security Scan Workflow
    ├─ Lint code
    ├─ Run tests
    ├─ Security scan
    ├─ Build multi-platform image → Push to registry
    ├─ Build single-platform image → Export as artifact
    └─ Upload artifact (available for 1 day)
         ↓
    (workflow_run trigger)
         ↓
Kubernetes Deployment Pipeline
    ├─ Helm chart tests
    ├─ Download Docker image artifact ← (Reuse from previous workflow!)
    ├─ Test deployment in Kind cluster
    └─ Deploy to staging/production
```

## Manual Triggering

You can manually trigger the Kubernetes deployment pipeline:

1. Go to GitHub Actions tab
2. Select "Kubernetes Deployment Pipeline"
3. Click "Run workflow"
4. Select the branch
5. Click "Run workflow"

This will run the K8s deployment using the most recent Docker image artifact from the last successful build on that branch.

## Important Notes

⚠️ **Artifact Retention**: The Docker image artifact is retained for 1 day. If you need to run the K8s deployment after that period, re-run the Docker build workflow first.

⚠️ **Branch Consistency**: The K8s deployment workflow runs on the same branch as the Docker build workflow to ensure it downloads the correct image artifact.

⚠️ **Conditional Execution**: The K8s deployment only runs if the Docker build workflow completed successfully.

## Testing the Setup

1. **Push to main branch:**
   ```bash
   git add .
   git commit -m "Update workflows"
   git push origin main
   ```

2. **Watch the workflows:**
   - First: "Docker Build and Security Scan" will run
   - Then: "Kubernetes Deployment Pipeline" will automatically start

3. **Check artifacts:**
   - Go to the completed Docker build workflow
   - Check "Artifacts" section
   - You should see "docker-image" artifact

## Troubleshooting

### If K8s deployment doesn't start automatically:
- Check that the Docker build workflow completed successfully
- Verify the branch names match in both workflows
- Check workflow permissions in repository settings

### If artifact download fails:
- Verify the artifact name matches: `docker-image`
- Check artifact retention period (default: 1 day)
- Ensure the workflow has permission to access artifacts

### If image loading fails in Kind:
- Verify the tar file was downloaded correctly
- Check Docker is running in the runner
- Ensure Kind cluster was created successfully
