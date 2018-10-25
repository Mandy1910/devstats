select
  c.sha,
  c.event_id,
  c.author_name,
  c.message,
  c.dup_actor_login,
  c.dup_repo_name,
  c.dup_created_at,
  c.encrypted_email,
  c.author_email,
  c.committer_name,
  c.committer_email,
  coalesce(c.dup_author_login, author.login, ''),
  coalesce(c.dup_committer_login, committer.login, ''),
  coalesce(r.org_login, ''),
  r.repo_group,
  r.alias,
  coalesce(actor.name, ''),
  coalesce(actor.country_id, ''),
  coalesce(actor.sex, ''),
  actor.sex_prob,
  coalesce(actor.tz, ''),
  actor.tz_offset,
  coalesce(actor.country_name, ''),
  coalesce(author.country_id, ''),
  coalesce(author.sex, ''),
  author.sex_prob,
  coalesce(author.tz, ''),
  author.tz_offset,
  coalesce(author.country_name, ''),
  coalesce(committer.country_id, ''),
  coalesce(committer.sex, ''),
  committer.sex_prob,
  coalesce(committer.tz, ''),
  committer.tz_offset,
  coalesce(committer.country_name, '')
from
  gha_commits c
left join
  gha_repos r
on
  c.dup_repo_id = r.id
  and c.dup_repo_name = r.name
left join
  gha_actors actor
on
  c.dup_actor_id = actor.id
left join
  gha_actors author
on
  c.author_id = author.id
left join
  gha_actors committer
on
  c.committer_id = committer.id
where
  dup_created_at >= '{{from}}'
  and dup_created_at < '{{to}}'
;
