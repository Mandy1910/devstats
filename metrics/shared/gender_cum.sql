select
  concat(inn.type, ';', case inn.sex when 'm' then 'Male' when 'f' then 'Female' end, '`', inn.repo_group, ';contributors,contributions,users,events') as name,
  inn.contributors,
  inn.contributions,
  inn.users,
  inn.events
from (
  select 'sexcum' as type,
    a.sex,
    'all' as repo_group,
    count(distinct e.actor_id) filter (where e.type in ('IssuesEvent', 'PullRequestEvent', 'PushEvent', 'CommitCommentEvent', 'IssueCommentEvent', 'PullRequestReviewCommentEvent')) as contributors,
    count(distinct e.id) filter (where e.type in ('IssuesEvent', 'PullRequestEvent', 'PushEvent', 'CommitCommentEvent', 'IssueCommentEvent', 'PullRequestReviewCommentEvent')) as contributions,
    count(distinct e.actor_id) as users,
    count(distinct e.id) as events
  from
    gha_events e,
    gha_actors a
  where
    (lower(e.dup_actor_login) {{exclude_bots}})
    and a.id = e.actor_id
    and a.sex is not null
    and a.sex != ''
    and a.sex_prob >= 0.7
    and e.created_at < '{{to}}'
  group by
    a.sex
  union select 'sexcum' as type,
    a.sex,
    coalesce(ecf.repo_group, r.repo_group) as repo_group,
    count(distinct e.actor_id) filter (where e.type in ('IssuesEvent', 'PullRequestEvent', 'PushEvent', 'CommitCommentEvent', 'IssueCommentEvent', 'PullRequestReviewCommentEvent')) as contributors,
    count(distinct e.id) filter (where e.type in ('IssuesEvent', 'PullRequestEvent', 'PushEvent', 'CommitCommentEvent', 'IssueCommentEvent', 'PullRequestReviewCommentEvent')) as contributions,
    count(distinct e.actor_id) as users,
    count(distinct e.id) as events
  from
    gha_repos r,
    gha_actors a,
    gha_events e
  left join
    gha_events_commits_files ecf
  on
    ecf.event_id = e.id
  where
    r.id = e.repo_id
    and (lower(e.dup_actor_login) {{exclude_bots}})
    and a.id = e.actor_id
    and a.sex is not null
    and a.sex != ''
    and a.sex_prob >= 0.7
    and e.created_at < '{{to}}'
  group by
    a.sex,
    coalesce(ecf.repo_group, r.repo_group)
) inn
where
  inn.repo_group is not null 
order by
  name
;
