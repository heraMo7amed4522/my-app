import { Module } from '@nestjs/common';
import { UserProgressService } from './user-progress.service';
import { UserProgressResolver } from './user-progress.resolver';

@Module({
  providers: [UserProgressService, UserProgressResolver],
})
export class UserProgressModule {}