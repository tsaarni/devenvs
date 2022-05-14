#!/usr/bin/env python3
#
# Pre-condition: install selenium-python and chrome webdriver
#
#   python3 -m venv selenium
#   . selenium/bin/activate
#   pip install selenium
#   pip install webdrivermanager
#   webdrivermanager chrome
#

from selenium.webdriver import Chrome

driver = Chrome(desired_capabilities={'acceptInsecureCerts': True})
driver.get("https://protected-oauth.127-0-0-101.nip.io/allowed")
assert "Sign in" in driver.title

driver.find_element_by_name("username").send_keys('joe')
driver.find_element_by_name("password").send_keys('joe')
driver.find_element_by_name("login").submit()

assert "OAuth flow failed" not in driver.page_source

#driver.close()
