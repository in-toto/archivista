import { Step } from "./Policy";
import { Token } from "./Auth";

export type Config = {
  image?: string;
  stages: string[];
} & Record<string, Job>;

export interface Job {
  image?: string;
  script: string[];
  stage: string;
  variables?: Record<string, string | number>;
  artifacts?: Artifacts;
}

export interface Artifacts {
  paths: string[];
}

function isJob(j: unknown): j is Job {
  return typeof j === "object" && j !== null && "script" in j;
}

/**
 * don't judge me... i wrote this at 6am after an all-nighter.
 * I wanted the steps coming out from this function to maintain the job order
 * by their stage in the gitlab config
 **/
export function StepsFromGitlabConfig(config: Config): Step[] {
  const jobsByStage = new Map<string, Map<string, Job>>();
  for (const prop in config) {
    const val = config[prop];
    if (isJob(val)) {
      let jobs = jobsByStage.get(val.stage);
      if (jobs === undefined) {
        jobs = new Map<string, Job>();
      }

      jobs = jobs.set(prop, val);
      jobsByStage.set(val.stage, jobs);
    }
  }

  const steps: Step[] = [];
  for (const stage of config.stages) {
    const stageJobs = jobsByStage.get(stage);
    if (stageJobs === undefined) {
      continue;
    }

    for (const [name, job] of stageJobs) {
      steps.push(jobToStep(name, job));
    }
  }

  return steps;
}

/**
 * We'll eventually want to use j to fill out extra information on the step..
 * but for now just punting
 **/
function jobToStep(name: string, _: Job): Step {
  return {
    name: name,
    functionaries: [],
    attestations: [{ predicate: "https://witness.testifysec.com/attestations/gitlab/v0.1", policies: [] }],
  };
}

export interface Project {
  id: number;
  name: string;
  name_with_namespace: string;
}

export interface User {
  name: string;
}

export interface Pipeline {
  id: number;
  project_id: number;
  ref: string;
  status: string;
  web_url: string;
}

export class ApiService {
  readonly #baseUrl: string;
  #token: Token;

  constructor(baseUrl: string, token: Token) {
    this.#baseUrl = baseUrl;
    this.#token = token;
  }

  async fetchProjects(): Promise<Project[]> {
    const url = new URL("/api/v4/projects", this.#baseUrl);
    url.searchParams.set("membership", "true");
    url.searchParams.set("order_by", "last_activity_at");

    const response = await fetch(url.toString(), {
      headers: {
        Authorization: `Bearer ${this.#token.access_token}`,
      },
    });

    if (response.status > 300) {
      return [] as Project[];
    }
    return (await response.json()) as Project[];
  }

  async fetchCurrentUser(): Promise<User> {
    const url = new URL("/oauth/userinfo", this.#baseUrl);
    const response = await fetch(url.toString(), {
      headers: {
        Authorization: `Bearer ${this.#token.access_token}`,
      },
    });

    return (await response.json()) as User;
  }

  async fetchPipelines(projectId: number): Promise<Pipeline[]> {
    const url = new URL(`/api/v4/projects/${projectId}/pipelines`, this.#baseUrl);
    const response = await fetch(url.toString(), {
      headers: {
        Authorization: `Bearer ${this.#token.access_token}`,
      },
    });

    return (await response.json()) as Pipeline[];
  }
}
