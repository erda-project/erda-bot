package gh

import (
	"context"

	"github.com/google/go-github/v35/github"
)

var client *github.Client

func MergePR(owner, repo string, pullNumber int) (*github.PullRequestMergeResult, error) {
	result, _, err := client.PullRequests.Merge(context.Background(), owner, repo, pullNumber, "", &github.PullRequestOptions{
		MergeMethod:        "squash",
		DontDefaultIfBlank: false,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ApprovePR(owner, repo string, pullNumber int) error {
	ctx := context.Background()
	_, _, err := client.PullRequests.CreateReview(ctx, owner, repo, pullNumber, &github.PullRequestReviewRequest{
		Event: &[]string{"APPROVE"}[0],
	})
	return err
}

func UpdateBranch(owner, repo string, pullNumber int) error {
	ctx := context.Background()
	_, _, err := client.PullRequests.UpdateBranch(ctx, owner, repo, pullNumber, nil)
	return err
}

func GetPR(owner, repo string, pullNumber int) (*github.PullRequest, error) {
	pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, pullNumber)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func GetBranchProtection(owner, repo, branch string) (*github.Protection, error) {
	protection, _, err := client.Repositories.GetBranchProtection(context.Background(), owner, repo, branch)
	if err != nil {
		return nil, err
	}
	return protection, nil
}

func AddIssueAssignees(owner, repo string, pullNumber int, assignees []string) error {
	_, _, err := client.Issues.AddAssignees(context.Background(), owner, repo, pullNumber, assignees)
	return err
}

func AddPRReviewers(owner, repo string, pullNumber int, reviewers []string) error {
	_, _, err := client.PullRequests.RequestReviewers(context.Background(), owner, repo, pullNumber, github.ReviewersRequest{Reviewers: reviewers})
	return err
}

func GetRepoStar(owner, repo string) (int, error) {
	r, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return 0, err
	}
	return r.GetStargazersCount(), nil
}
