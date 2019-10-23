import { TestBed } from '@angular/core/testing';

import { RoleBindingService } from './role-binding.service';

describe('RoleBindingService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: RoleBindingService = TestBed.get(RoleBindingService);
    expect(service).toBeTruthy();
  });
});
