export class RoleBinding {
  name: string;
  createdAt: Date;
  subjects: Subject[];
  roleRef: string;

}

export interface Subject {
  type: string;
  value: string;
}
