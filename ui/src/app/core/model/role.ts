export class Role {
  name: string;
  displayName: string;
  createdAt: Date;
  rules: Rule[];
}

export class Rule {
  // TODO: change model to 'resource: string' only
  resources: string[];
  actions: string[];
}
