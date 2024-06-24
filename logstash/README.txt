
# Setup logstash environment, generic to all plugin development

cd ~/work/logstash
rbenv install  # installs version defined in .ruby-version

# activate ruby version
rbenv global jruby-9.3.4.0
ruby --version



# install some dependencies for jruby
gem install rake
gem install bundler

rake bootstrap                   # install jruby and dependencies
rake plugin:install-default      # install default plugins under vendor/bundle


bin/bundle show --paths          # show paths where gems are loaded from


# if getting error
#   Gem::GemNotFoundException: can't find gem rake (>= 0.a) with executable rake
#
# or
#
# Execution failed for task ':bootstrap'.
# > (VersionConflict) Bundler could not find compatible versions for gem "aws-sdk-core":
#    In Gemfile:
#      logstash-output-sns was resolved to 4.0.8, which depends on
#        logstash-mixin-aws (>= 1.0.0) was resolved to 4.0.2, which depends on
#
# then just clean the repo from modified and untracked files and start from the beginning

git clean -fdx                     # remove untracked files
git reset --hard                   # clean modified files (working tree + index)




### Run java tests

java --version
# openjdk 17.0.10 2024-01-16
# OpenJDK Runtime Environment Temurin-17.0.10+7 (build 17.0.10+7)
# OpenJDK 64-Bit Server VM Temurin-17.0.10+7 (build 17.0.10+7, mixed mode, sharing)

tar zxf ~/Downloads/logstash-8.12.2.tar.gz
cd logstash-8.12.2/
rbenv install
rbenv global jruby-9.3.10.0
ruby --version
# jruby 9.3.10.0 (2.6.8) 2023-02-01 107b2e6697 OpenJDK 64-Bit Server VM 17.0.10+7 on 17.0.10+7 +jit [x86_64-linux]

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=$PWD

./gradlew installDevelopmentGems
./gradlew installDefaultGems
./gradlew javaTests
