import {Component, Input, OnInit} from '@angular/core';
import {Role} from "../core/model/role";
import {RoleService} from "../core/services/role.service";
import {ActivatedRoute} from "@angular/router";

interface RuleDetail {
  resource: string;
  actions: Map<string, boolean>
}

@Component({
  selector: 'app-role-detail',
  templateUrl: './role-detail.component.html',
  styleUrls: ['./role-detail.component.scss']
})
export class RoleDetailComponent implements OnInit {
  displayedColumns: string[] = ['resource', 'list', 'get', 'create', 'delete', 'update', 'publish', 'unpublish'];
  private ruleDetails: RuleDetail[] = [];

  @Input() role: Role;

  constructor(private route: ActivatedRoute,
              private roleService: RoleService) {
  }

  ngOnInit() {
    this.getRole();
  }

  private getRole() {
    const name = this.route.snapshot.paramMap.get("name");
    this.roleService.getRole(name)
      .subscribe(
        role => this.updateDataSource(role),
        err => console.error(`failed to get role detail of ${name}: ${err}`)
      );
  }

  private updateDataSource(role: Role) {
    this.role = role;
    for (let r of role.rules) {
      let ruleDetail: RuleDetail = {
        resource: r.resources[0],
        actions: new Map<string, boolean>()
      };
      ruleDetail.actions.clear();
      for (let a of r.actions) {
        ruleDetail.actions['x' + a] = true
      }
      this.ruleDetails.push(ruleDetail)
    }
    console.log(`Role: ${JSON.stringify(this.role)}`);
    console.log(this.ruleDetails)
  }

  private getActionValue(rule: RuleDetail, action: string): string {
    let value = rule.actions['x' + action];
    console.log(`${action}=${value}`);
    return value ? value : false;
  }
}
