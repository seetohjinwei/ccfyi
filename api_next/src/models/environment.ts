export interface Environment {
  workspace: string;
  environment: string;
  environment_variables: EnvironmentVariable[];
}

type EnvironmentVariableType = "default" | "secret";

export interface EnvironmentVariable {
  workspace: string;
  environment: string;
  variable: string;
  type: EnvironmentVariableType;
  initial_value: string;
  current_value: string;
}
