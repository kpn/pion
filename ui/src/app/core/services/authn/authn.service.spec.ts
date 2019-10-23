import { TestBed } from '@angular/core/testing';

import { AuthnService } from './authn.service';

describe('AuthnService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: AuthnService = TestBed.get(AuthnService);
    expect(service).toBeTruthy();
  });
});
