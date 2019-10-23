import { TestBed, async, inject } from '@angular/core/testing';

import { AuthnGuard } from './authn.guard';

describe('AuthnGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AuthnGuard]
    });
  });

  it('should ...', inject([AuthnGuard], (guard: AuthnGuard) => {
    expect(guard).toBeTruthy();
  }));
});
