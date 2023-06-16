#!/usr/bin/env python3
#
# Pre-condition: install selenium-python and chrome webdriver
#
#   python3 -m venv .venv
#   . .venv/bin/activate
#   pip install selenium webdrivermanager
#   webdrivermanager chrome
#

from selenium.webdriver import Chrome
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC


try:

    driver = Chrome()
    driver.get("http://localhost:8080/realms/master/account/")
    assert "Keycloak Account Management" in driver.title

    sign_in = WebDriverWait(driver, 10).until(
        EC.element_to_be_clickable((By.ID, "landingSignInButton"))
    )

    sign_in.click()

    driver.find_element(By.NAME, "username").send_keys('mustchange')
    driver.find_element(By.NAME, "password").send_keys('mustchange')
    driver.find_element(By.NAME, "login").submit()

    assert "Update password" in driver.page_source, "Update password page not found"

except Exception as e:
    print(f"An exception occurred: {e}")

finally:
    input("Press Enter to continue...")
    driver.quit()
