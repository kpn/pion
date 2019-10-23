import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RoleBindingsComponent } from './role-bindings.component';

describe('RoleBindingsComponent', () => {
  let component: RoleBindingsComponent;
  let fixture: ComponentFixture<RoleBindingsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RoleBindingsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RoleBindingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
