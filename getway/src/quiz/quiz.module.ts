import { Module } from '@nestjs/common';
import { QuizService } from './quiz.service';
import { QuizResolver } from './quiz.resolver';

@Module({
  providers: [QuizService, QuizResolver],
  exports: [QuizService],
})
export class QuizModule {}