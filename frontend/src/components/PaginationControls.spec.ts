import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import PaginationControls from './PaginationControls.vue'

function buttonLabels(wrapper: ReturnType<typeof mount>) {
  return wrapper.findAll('.page-btn').map((b) => b.text())
}

describe('PaginationControls', () => {
  it('renders nothing when there is a single page', () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 0, totalPages: 1 },
    })
    expect(wrapper.find('nav').exists()).toBe(false)
  })

  it('lists every page (1-indexed labels) when total <= 7', () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 0, totalPages: 5 },
    })
    // prev arrow + 5 numbered pages + next arrow
    expect(buttonLabels(wrapper)).toEqual(['←', '1', '2', '3', '4', '5', '→'])
    expect(wrapper.find('.page-ellipsis').exists()).toBe(false)
  })

  it('shows a trailing ellipsis near the start of a long range', () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 0, totalPages: 10 },
    })
    // first, current window, ellipsis, last
    expect(buttonLabels(wrapper)).toEqual(['←', '1', '2', '10', '→'])
    expect(wrapper.findAll('.page-ellipsis')).toHaveLength(1)
  })

  it('shows both ellipses when the current page is in the middle', () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 5, totalPages: 12 },
    })
    expect(buttonLabels(wrapper)).toEqual(['←', '1', '5', '6', '7', '12', '→'])
    expect(wrapper.findAll('.page-ellipsis')).toHaveLength(2)
  })

  it('marks the current page active', () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 2, totalPages: 5 },
    })
    const active = wrapper.find('.page-btn.active')
    expect(active.text()).toBe('3')
  })

  it('disables the prev arrow on the first page and next arrow on the last', () => {
    const first = mount(PaginationControls, { props: { currentPage: 0, totalPages: 5 } })
    const buttons = first.findAll('.page-btn')
    expect(buttons[0].attributes('disabled')).toBeDefined() // prev
    expect(buttons[buttons.length - 1].attributes('disabled')).toBeUndefined() // next

    const last = mount(PaginationControls, { props: { currentPage: 4, totalPages: 5 } })
    const lastButtons = last.findAll('.page-btn')
    expect(lastButtons[0].attributes('disabled')).toBeUndefined()
    expect(lastButtons[lastButtons.length - 1].attributes('disabled')).toBeDefined()
  })

  it('emits update:currentPage with the zero-based index on click', async () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 0, totalPages: 5 },
    })
    // Click the button labelled "3" -> zero-based page 2
    const target = wrapper.findAll('.page-btn').find((b) => b.text() === '3')!
    await target.trigger('click')
    expect(wrapper.emitted('update:currentPage')).toEqual([[2]])
  })

  it('does not emit when navigating out of bounds', async () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 0, totalPages: 5 },
    })
    // prev arrow while on first page -> goTo(-1), rejected by guard
    await wrapper.findAll('.page-btn')[0].trigger('click')
    expect(wrapper.emitted('update:currentPage')).toBeUndefined()
  })
})
