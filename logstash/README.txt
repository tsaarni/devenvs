
# Setup logstash environment, generic to all plugin development

cd ~/work/logstash
rbenv install  # installs version defined in logstash/.ruby-version

# activate ruby version
rbenv global jruby-9.3.4.0
ruby --version




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
