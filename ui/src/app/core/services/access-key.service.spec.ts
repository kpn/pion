import { TestBed } from '@angular/core/testing';

import { AccessKeyService } from './access-key.service';

describe('AccessKeyService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: AccessKeyService = TestBed.get(AccessKeyService);
    expect(service).toBeTruthy();
  });
});
